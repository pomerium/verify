package verify

import (
	"bytes"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"html/template"
	"io"
	"io/fs"
	stdlog "log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/rs/zerolog/log"

	sdk "github.com/pomerium/sdk-go"
)

//go:embed ui/dist
var uiFS embed.FS

func (srv *Server) initRouter() {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialTLSContext = srv.tlsVerifier.DialTLSContext
	client := &http.Client{
		Transport: transport,
		Timeout:   maxRemoteWait,
	}

	expected := &jwt.Expected{
		Issuer: srv.cfg.expectedJWTIssuer,
	}
	if aud := srv.cfg.expectedJWTAudience; aud != "" {
		expected.Audience = jwt.Audience([]string{aud})
	}

	datastore, err := sdk.NewLRUKeyStore(1024)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	verifier, err := sdk.New(&sdk.Options{
		Datastore:    datastore,
		HTTPClient:   client,
		Logger:       stdlog.New(log.With().Logger(), "", 0),
		JWKSEndpoint: srv.cfg.jwksEndpoint,
		Expected:     expected,
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	srv.router = chi.NewRouter()
	srv.router.Use(sdk.AddIdentityToRequest(verifier))

	srv.router.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// mount api
	srv.router.Route("/api", func(r chi.Router) {
		r.Use(middleware.NoCache)

		r.Get("/verify-info", srv.serveAPIVerifyInfo)
		r.Post("/webauthn-authenticate", srv.serveAPIWebAuthnAuthenticate)
		r.Post("/webauthn-register", srv.serveAPIWebAuthnRegister)

		r.NotFound(func(w http.ResponseWriter, r *http.Request) {
			http.NotFound(w, r)
		})
	})
	srv.router.Get("/headers", srv.serveHeaders)
	srv.router.Get("/json", srv.serveAPIVerifyInfo)

	// mount static files
	root := "ui/dist"
	etags := map[string]string{}
	_ = fs.WalkDir(uiFS, root, func(p string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		etags[p] = computeEtag(p)

		srv.router.Get(p[len(root):], func(w http.ResponseWriter, r *http.Request) {
			srv.serveStatic(w, r, p, etags[p])
		})

		return nil
	})
	srv.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if path.Ext(r.URL.Path) == "" {
			p := path.Join(root, "index.html")
			srv.serveStatic(w, r, p, etags[p])
			return
		}

		http.NotFound(w, r)
	})
}

func (srv *Server) serveHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(r.Header)
}

func (srv *Server) serveStatic(w http.ResponseWriter, r *http.Request, name, etag string) {
	if strings.HasSuffix(name, ".html") {
		srv.serveTemplate(w, r, name, etag)
		return
	}

	f, err := uiFS.Open(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	rs, ok := f.(io.ReadSeeker)
	if !ok {
		log.Fatal().Str("name", name).Msg("invalid static file, not a read seeker")
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", `"`+etag+`"`)
	http.ServeContent(w, r, path.Base(name), time.Time{}, rs)
}

func (srv *Server) serveTemplate(w http.ResponseWriter, r *http.Request, name, etag string) {
	bs, err := uiFS.ReadFile(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tpl, err := template.New("").Parse(string(bs))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{}
	if v, ok := os.LookupEnv("GOOGLE_TAG_MANAGER_ID"); ok {
		data["GoogleTagManagerID"] = v
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Etag", `"`+etag+`"`)
	http.ServeContent(w, r, path.Base(name), time.Time{}, bytes.NewReader(buf.Bytes()))
}

func (srv *Server) serveAPIVerifyInfo(w http.ResponseWriter, r *http.Request) {
	type M = map[string]interface{}

	res := M{
		"headers": getPomeriumHeaders(r),
	}
	var tlsErrStr string
	if identity, err := sdk.FromContext(r.Context()); err == nil {
		res["identity"] = identity
		jwksDomain := identity.Issuer
		if srv.cfg.jwksEndpoint != "" {
			u, err := url.Parse(srv.cfg.jwksEndpoint)
			if err == nil {
				jwksDomain = u.Hostname()
			}
		}
		if e := srv.tlsVerifier.GetTLSError(jwksDomain); e != nil {
			tlsErrStr = e.Error()
		}
	} else {
		res["identity"] = getUnverifiedIdentity(r)
		res["error"] = err.Error()
	}
	res["request"] = M{
		"origin":   getOrigin(r),
		"method":   r.Method,
		"url":      r.URL.RequestURI(),
		"host":     r.Host,
		"hostname": getHostname(),
		"tlsError": tlsErrStr,
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

func getUnverifiedIdentity(r *http.Request) *sdk.Identity {
	identity := new(sdk.Identity)
	jwt, err := jose.ParseSigned(sdk.TokenFromHeader(r))
	if err != nil {
		log.Err(err).Msg("error parsing JWT assertion header")
		return identity
	}

	err = json.Unmarshal(jwt.UnsafePayloadWithoutVerification(), identity)
	if err != nil {
		log.Err(err).Msg("error unmarshaling JWT assertion header")
		return identity
	}

	return identity
}

func getHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func getOrigin(r *http.Request) string {
	origin := r.Header.Get("X-Forwarded-For")
	if origin == "" {
		origin = r.RemoteAddr
	}
	return origin
}

func getPomeriumHeaders(r *http.Request) http.Header {
	hdrs := r.Header.Clone()
	for v := range r.Header {
		v = strings.ToLower(v)
		if !strings.Contains(v, "x-pomerium-claim") {
			hdrs.Del(v)
		}
	}
	return hdrs
}

func computeEtag(name string) string {
	bs, _ := uiFS.ReadFile(name)
	h := sha256.Sum256(bs)
	return hex.EncodeToString(h[:])
}

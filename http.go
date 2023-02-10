package verify

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	stdlog "log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"

	sdk "github.com/pomerium/sdk-go"
)

//go:embed ui/dist
var uiFS embed.FS

func (srv *Server) initRouter(jwksEndpoint string) {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.DialTLSContext = srv.tlsVerifier.DialTLSContext
	client := &http.Client{
		Transport: transport,
	}

	verifier, err := sdk.New(&sdk.Options{
		Datastore:    NewCache(1024),
		HTTPClient:   client,
		Logger:       stdlog.New(log.With().Logger(), "", 0),
		JWKSEndpoint: jwksEndpoint,
	})
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	srv.router = chi.NewRouter()
	srv.router.Use(sdk.AddIdentityToRequest(verifier))

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
	json.NewEncoder(w).Encode(r.Header)
}

func (srv *Server) serveStatic(w http.ResponseWriter, r *http.Request, name, etag string) {
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

func (srv *Server) serveAPIVerifyInfo(w http.ResponseWriter, r *http.Request) {
	type M = map[string]interface{}

	res := M{
		"headers": getPomeriumHeaders(r),
	}
	var tlsErrStr string
	if identity, err := sdk.FromContext(r.Context()); err == nil {
		res["identity"] = identity
		if e := srv.tlsVerifier.GetTLSError(identity.Issuer); e != nil {
			tlsErrStr = e.Error()
		}
	} else {
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
	json.NewEncoder(w).Encode(res)
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

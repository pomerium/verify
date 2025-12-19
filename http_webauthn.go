package verify

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/pomerium/verify/internal/storage"
	"github.com/pomerium/webauthn"
)

const maxBodySize = 4 * 1024 * 1024

func (srv *Server) serveAPIWebAuthnAuthenticate(w http.ResponseWriter, r *http.Request) {
	var req storage.WebAuthnAuthenticateRequest
	err := decodeJSONBody(r, &req)
	if err != nil {
		log.Error().Err(err).Msg("bad request for webauthn authenticate")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = srv.storage.SetAuthenticateRequest(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	origin := getRPOrigin(r)
	if clientData, err := req.Credential.Response.UnmarshalClientData(); err == nil {
		// use the provided origin
		origin = clientData.Origin
	}

	rp := webauthn.NewRelyingParty(origin, srv.storage)

	// This verification is insecure because we trust the options provided by the client. A real implementation
	// should generate the options (particularly the challenge) on the server.
	credential, err := rp.VerifyAuthenticationCeremony(r.Context(), req.Options, req.Credential)
	if err != nil {
		log.Error().Err(err).Msg("webauthn: invalid authentication ceremony")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	_ = credential
}

func (srv *Server) serveAPIWebAuthnRegister(w http.ResponseWriter, r *http.Request) {
	var req storage.WebAuthnRegisterRequest
	err := decodeJSONBody(r, &req)
	if err != nil {
		log.Error().Err(err).Msg("bad request for webauthn register")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = srv.storage.SetRegisterRequest(r.Context(), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	origin := getRPOrigin(r)
	if clientData, err := req.Credential.Response.UnmarshalClientData(); err == nil {
		// use the provided origin
		origin = clientData.Origin
	}

	rp := webauthn.NewRelyingParty(origin, srv.storage)

	// This verification is insecure because we trust the options provided by the client. A real implementation
	// should generate the options (particularly the challenge) on the server.
	credential, err := rp.VerifyRegistrationCeremony(r.Context(), req.Options, req.Credential)
	if err != nil {
		log.Error().Err(err).Msg("webauthn: invalid registration ceremony")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	err = srv.storage.SetCredential(r.Context(), credential)
	if err != nil {
		log.Error().Err(err).Msg("webauthn: invalid registration ceremony")
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}

func decodeJSONBody(r *http.Request, dsts ...interface{}) error {
	defer func() { _ = r.Body.Close() }()

	bs, err := io.ReadAll(io.LimitReader(r.Body, maxBodySize))
	if err != nil {
		return err
	}

	for _, dst := range dsts {
		err = json.Unmarshal(bs, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func getRPOrigin(r *http.Request) string {
	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}
	return scheme + "://" + r.Host
}

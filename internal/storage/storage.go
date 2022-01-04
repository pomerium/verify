package storage

import (
	"context"
	"encoding/base64"

	"github.com/pomerium/webauthn"
)

// A Backend stores credentials and requests.
type Backend interface {
	GetCredential(ctx context.Context, credentialID []byte) (*webauthn.Credential, error)
	SetAuthenticateRequest(ctx context.Context, req *WebAuthnAuthenticateRequest) error
	SetCredential(ctx context.Context, credential *webauthn.Credential) error
	SetRegisterRequest(ctx context.Context, req *WebAuthnRegisterRequest) error
}

// WebAuthnAuthenticateRequest is the authenticate request info and credential.
type WebAuthnAuthenticateRequest struct {
	Options    *webauthn.PublicKeyCredentialRequestOptions `json:"options"`
	Credential *webauthn.PublicKeyAssertionCredential      `json:"credential"`
}

// GetID gets the ID for the request.
func (req WebAuthnAuthenticateRequest) GetID() string {
	return base64.RawURLEncoding.EncodeToString(req.Credential.RawID)
}

// WebAuthnRegisterRequest is the register request info and credential.
type WebAuthnRegisterRequest struct {
	Options    *webauthn.PublicKeyCredentialCreationOptions `json:"options"`
	Credential *webauthn.PublicKeyCreationCredential        `json:"credential"`
}

// GetID gets the ID for the request.
func (req WebAuthnRegisterRequest) GetID() string {
	return base64.RawURLEncoding.EncodeToString(req.Credential.RawID)
}

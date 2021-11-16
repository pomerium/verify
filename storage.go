package verify

import (
	"context"
	"encoding/base64"

	"cloud.google.com/go/firestore"
	"github.com/pomerium/webauthn"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collectionNameWebAuthnAuthenticateRequests = "webauthn-authenticate-requests"
	collectionNameWebAuthnCredentials          = "webauthn-credentials"
	collectionNameWebAuthnRegisterRequests     = "webauthn-register-requests"
)

// Storage is used to store data in firestore.
type Storage struct {
	client *firestore.Client
}

// GetCredential gets a WebAuthn credential from storage.
func (s *Storage) GetCredential(ctx context.Context, credentialID []byte) (*webauthn.Credential, error) {
	id := base64.RawURLEncoding.EncodeToString(credentialID)
	var credential webauthn.Credential
	err := s.get(ctx, collectionNameWebAuthnCredentials, id, &credential)
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

// SetAuthenticateRequest sets a WebAuthn authenticate request in storage.
func (s *Storage) SetAuthenticateRequest(ctx context.Context, req *WebAuthnAuthenticateRequest) error {
	return s.set(ctx, collectionNameWebAuthnAuthenticateRequests, req.GetID(), req)
}

// SetCredential sets a WebAuthn credential in storage.
func (s *Storage) SetCredential(ctx context.Context, credential *webauthn.Credential) error {
	id := base64.RawURLEncoding.EncodeToString(credential.ID)
	return s.set(ctx, collectionNameWebAuthnCredentials, id, credential)
}

// SetRegisterRequest sets a WebAuthn register request in storage.
func (s *Storage) SetRegisterRequest(ctx context.Context, req *WebAuthnRegisterRequest) error {
	return s.set(ctx, collectionNameWebAuthnRegisterRequests, req.GetID(), req)
}

func (s *Storage) get(ctx context.Context, collectionName, objectID string, dst interface{}) error {
	collection := s.client.Collection(collectionName)
	doc := collection.Doc(objectID)

	snapshot, err := doc.Get(ctx)
	if status.Code(err) == codes.NotFound {
		return webauthn.ErrCredentialNotFound
	} else if err != nil {
		log.Error().Err(err).Msg("storage: failed to get object")
		return err
	}

	log.Info().Str("id", doc.ID).Str("path", doc.Path).Msg("storage: got object")
	return snapshot.DataTo(dst)
}

func (s *Storage) set(ctx context.Context, collectionName, objectID string, obj interface{}) error {
	collection := s.client.Collection(collectionName)
	doc := collection.Doc(objectID)
	_, err := doc.Set(ctx, obj)
	if err != nil {
		log.Error().Err(err).Msg("storage: failed to set object")
		return err
	}
	log.Info().Str("id", doc.ID).Str("path", doc.Path).Msg("storage: set object")
	return nil
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

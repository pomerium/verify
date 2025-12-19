package storage

import (
	"context"
	"encoding/base64"

	"cloud.google.com/go/firestore"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pomerium/webauthn"
)

const (
	collectionNameWebAuthnAuthenticateRequests = "webauthn-authenticate-requests"
	collectionNameWebAuthnCredentials          = "webauthn-credentials"
	collectionNameWebAuthnRegisterRequests     = "webauthn-register-requests"
)

// FirestoreBackend is used to store data in firestore.
type FirestoreBackend struct {
	client *firestore.Client
}

// NewFirestoreBackend creates a new Firestore backend.
func NewFirestoreBackend(client *firestore.Client) *FirestoreBackend {
	return &FirestoreBackend{
		client: client,
	}
}

// GetCredential gets a WebAuthn credential from storage.
func (backend *FirestoreBackend) GetCredential(ctx context.Context, credentialID []byte) (*webauthn.Credential, error) {
	id := base64.RawURLEncoding.EncodeToString(credentialID)
	var credential webauthn.Credential
	err := backend.get(ctx, collectionNameWebAuthnCredentials, id, &credential)
	if err != nil {
		return nil, err
	}
	return &credential, nil
}

// SetAuthenticateRequest sets a WebAuthn authenticate request in storage.
func (backend *FirestoreBackend) SetAuthenticateRequest(ctx context.Context, req *WebAuthnAuthenticateRequest) error {
	return backend.set(ctx, collectionNameWebAuthnAuthenticateRequests, req.GetID(), req)
}

// SetCredential sets a WebAuthn credential in storage.
func (backend *FirestoreBackend) SetCredential(ctx context.Context, credential *webauthn.Credential) error {
	id := base64.RawURLEncoding.EncodeToString(credential.ID)
	return backend.set(ctx, collectionNameWebAuthnCredentials, id, credential)
}

// SetRegisterRequest sets a WebAuthn register request in storage.
func (backend *FirestoreBackend) SetRegisterRequest(ctx context.Context, req *WebAuthnRegisterRequest) error {
	return backend.set(ctx, collectionNameWebAuthnRegisterRequests, req.GetID(), req)
}

func (backend *FirestoreBackend) get(ctx context.Context, collectionName, objectID string, dst interface{}) error {
	collection := backend.client.Collection(collectionName)
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

func (backend *FirestoreBackend) set(ctx context.Context, collectionName, objectID string, obj interface{}) error {
	collection := backend.client.Collection(collectionName)
	doc := collection.Doc(objectID)
	_, err := doc.Set(ctx, obj)
	if err != nil {
		log.Error().Err(err).Msg("storage: failed to set object")
		return err
	}
	log.Info().Str("id", doc.ID).Str("path", doc.Path).Msg("storage: set object")
	return nil
}

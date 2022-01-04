package storage

import (
	"context"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/pomerium/webauthn"
)

// A InMemoryBackend stores data in-memory.
type InMemoryBackend struct {
	mu          sync.Mutex
	credentials map[string]*webauthn.Credential
}

// NewInMemoryBackend creates a new InMemoryBackend.
func NewInMemoryBackend() *InMemoryBackend {
	return &InMemoryBackend{}
}

// GetCredential retrieves a credential.
func (backend *InMemoryBackend) GetCredential(ctx context.Context, credentialID []byte) (*webauthn.Credential, error) {
	key := base64.RawURLEncoding.EncodeToString(credentialID)
	backend.mu.Lock()
	credential, ok := backend.credentials[key]
	backend.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return credential, nil
}

// SetAuthenticateRequest is a no-op.
func (backend *InMemoryBackend) SetAuthenticateRequest(ctx context.Context, req *WebAuthnAuthenticateRequest) error {
	// ignored
	return nil
}

// SetCredential saves a credential.
func (backend *InMemoryBackend) SetCredential(ctx context.Context, credential *webauthn.Credential) error {
	key := base64.RawURLEncoding.EncodeToString(credential.ID)
	backend.mu.Lock()
	backend.credentials[key] = credential
	backend.mu.Unlock()
	return nil
}

// SetRegisterRequest is a no-op.
func (backend *InMemoryBackend) SetRegisterRequest(ctx context.Context, req *WebAuthnRegisterRequest) error {
	// ignored
	return nil
}

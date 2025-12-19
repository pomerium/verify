package storage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/pomerium/webauthn"
)

func TestInMemoryBackend(t *testing.T) {
	ctx, clearTimeout := context.WithTimeout(context.Background(), time.Second*10)
	defer clearTimeout()

	backend := NewInMemoryBackend()
	err := backend.SetCredential(ctx, &webauthn.Credential{
		ID: []byte{1, 2, 3, 4},
	})
	assert.NoError(t, err)
	cred, err := backend.GetCredential(ctx, []byte{1, 2, 3, 4})
	assert.NoError(t, err)
	assert.NotNil(t, cred)
}

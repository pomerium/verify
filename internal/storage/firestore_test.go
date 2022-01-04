package storage

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pomerium/webauthn"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestFirestoreBackend(t *testing.T) {
	assert.NoError(t, WithTestFirestore(func(ctx context.Context, client *firestore.Client) error {
		backend := NewFirestoreBackend(client)
		err := backend.SetCredential(ctx, &webauthn.Credential{
			ID: []byte{1, 2, 3, 4},
		})
		assert.NoError(t, err)

		cred, err := backend.GetCredential(ctx, []byte{1, 2, 3, 4})
		assert.NoError(t, err)
		assert.NotNil(t, cred)

		return nil
	}))
}

// WithTestFirestore starts a test firestore emulator.
func WithTestFirestore(handler func(context.Context, *firestore.Client) error) error {
	maxWait := time.Minute * 10

	ctx, clearTimeout := context.WithTimeout(context.Background(), maxWait)
	defer clearTimeout()

	// if we've already set an emulator, don't start a new one
	if _, ok := os.LookupEnv("FIRESTORE_EMULATOR_HOST"); ok {
		client, err := firestore.NewClient(ctx, "test")
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
		defer client.Close()
		return handler(ctx, client)
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		return err
	}

	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "andreysenov/firebase-tools",
		Tag:        "latest",
		Cmd: []string{
			"firebase",
			"emulators:start",
			"--project",
			"test",
		},
		Mounts: []string{
			filepath.Join(wd, "..", "..") + ":/home/node",
		},
		ExposedPorts: []string{"8080"},
	})
	if err != nil {
		return err
	}
	_ = resource.Expire(uint(maxWait.Seconds()))
	go tailLogs(ctx, pool, resource)

	host := resource.Container.NetworkSettings.IPAddress + ":8080"

	defer func() { os.Unsetenv("FIRESTORE_EMULATOR_HOST") }()
	os.Setenv("FIRESTORE_EMULATOR_HOST", host)

	var client *firestore.Client
	if err := pool.Retry(func() error {
		log.Info().Str("host", host).Msg("connecting to firestore")
		var err error
		client, err = firestore.NewClient(ctx, "test")
		if err != nil {
			log.Error().Err(err).Send()
			return err
		}
		return nil
	}); err != nil {
		_ = pool.Purge(resource)
		return err
	}

	e := handler(ctx, client)
	_ = client.Close()

	if err := pool.Purge(resource); err != nil {
		return err
	}

	return e
}

func tailLogs(ctx context.Context, pool *dockertest.Pool, resource *dockertest.Resource) {
	_ = pool.Client.Logs(docker.LogsOptions{
		Context:      ctx,
		Stderr:       true,
		Stdout:       true,
		Follow:       true,
		Container:    resource.Container.ID,
		OutputStream: os.Stderr,
	})
}

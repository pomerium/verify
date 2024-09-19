package storage

import (
	"context"
	"net"
	"testing"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pomerium/webauthn"
)

func TestFirestoreBackend(t *testing.T) {
	ctx, clearTimeout := context.WithTimeout(context.Background(), time.Minute*10)
	defer clearTimeout()

	WithTestFirestore(t, func(client *firestore.Client) {
		backend := NewFirestoreBackend(client)
		err := backend.SetCredential(ctx, &webauthn.Credential{
			ID: []byte{1, 2, 3, 4},
		})
		assert.NoError(t, err)

		cred, err := backend.GetCredential(ctx, []byte{1, 2, 3, 4})
		assert.NoError(t, err)
		assert.NotNil(t, cred)
	})
}

func WithTestFirestore(t *testing.T, f func(client *firestore.Client)) {
	ctx, clearTimeout := context.WithTimeout(context.Background(), time.Minute*10)
	defer clearTimeout()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Name:         "pomerium-verify-firebase",
			Image:        "andreysenov/firebase-tools:latest",
			ExposedPorts: []string{"8080/tcp"},
			Cmd: []string{
				"firebase",
				"emulators:start",
				"--project",
				"test",
			},
			WaitingFor: wait.ForAll(
				wait.ForListeningPort("8080"),
				wait.ForLog("Emulator Hub running"),
			),
			Files: []testcontainers.ContainerFile{{
				HostFilePath:      "../../firebase.json",
				ContainerFilePath: "/home/node/firebase.json",
				FileMode:          0o644,
			}},
		},
		Started: true,
		Logger:  testcontainers.TestLogger(t),
		Reuse:   true,
	})
	require.NoError(t, err)

	mappedPort, err := container.MappedPort(ctx, "8080")
	require.NoError(t, err)

	t.Setenv("FIRESTORE_EMULATOR_HOST", net.JoinHostPort("127.0.0.1", mappedPort.Port()))
	client, err := firestore.NewClient(ctx, "test")
	require.NoError(t, err)

	f(client)
}

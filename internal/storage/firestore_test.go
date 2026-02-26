package storage

import (
	"net"
	"testing"

	"cloud.google.com/go/firestore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/log"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/pomerium/webauthn"
)

func TestFirestoreBackend(t *testing.T) {
	client := StartFirestore(t)

	backend := NewFirestoreBackend(client)
	err := backend.SetCredential(t.Context(), &webauthn.Credential{
		ID: []byte{1, 2, 3, 4},
	})
	assert.NoError(t, err)

	cred, err := backend.GetCredential(t.Context(), []byte{1, 2, 3, 4})
	assert.NoError(t, err)
	assert.NotNil(t, cred)
}

func StartFirestore(tb testing.TB) *firestore.Client {
	ctx := tb.Context()

	container, err := testcontainers.Run(ctx, "andreysenov/firebase-tools:13.32.0",
		testcontainers.WithLogger(log.TestLogger(tb)),
		testcontainers.WithExposedPorts("8080/tcp"),
		testcontainers.WithCmd("firebase", "emulators:start", "--project", "test"),
		testcontainers.WithWaitStrategy(
			wait.ForListeningPort("8080"),
			wait.ForLog("All emulators ready!"),
		),
		testcontainers.WithFiles(testcontainers.ContainerFile{
			HostFilePath:      "../../firebase.json",
			ContainerFilePath: "/home/node/firebase.json",
			FileMode:          0o644,
		}),
	)
	require.NoError(tb, err)

	host, err := container.Host(ctx)
	require.NoError(tb, err)
	port, err := container.MappedPort(ctx, "8080")
	require.NoError(tb, err)
	tb.Setenv("FIRESTORE_EMULATOR_HOST", net.JoinHostPort(host, port.Port()))

	client, err := firestore.NewClient(ctx, "test")
	require.NoError(tb, err)
	return client
}

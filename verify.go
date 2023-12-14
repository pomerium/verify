package verify

import (
	"context"
	"net"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/pomerium/verify/internal/storage"
)

// Server is the verify server backend.
type Server struct {
	cfg *config

	http        *http.Server
	router      chi.Router
	storage     storage.Backend
	tlsVerifier *tlsVerifier
}

// New creates a new Server.
func New(options ...Option) *Server {
	cfg := getConfig(options...)
	return &Server{
		cfg:         cfg,
		tlsVerifier: newTLSVerifier(),
	}
}

// Run runs the server.
func (srv *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	err := srv.init(ctx)
	if err != nil {
		return err
	}

	eg.Go(func() error {
		log.Info().
			Str("bind-addr", srv.cfg.bindAddress).
			Msg("starting http server")
		return srv.http.ListenAndServe()
	})
	return eg.Wait()
}

func (srv *Server) init(ctx context.Context) error {
	log.Info().
		Str("project-id", srv.cfg.firestoreProjectID).
		Msg("connecting to firestore")
	client, err := firestore.NewClient(ctx, srv.cfg.firestoreProjectID)
	if err == nil {
		srv.storage = storage.NewFirestoreBackend(client)
	} else {
		log.Error().Err(err).Msg("failed to create firestore client, falling back to in-memory storage")
		srv.storage = storage.NewInMemoryBackend()
	}

	srv.initRouter()

	srv.http = &http.Server{
		Addr: srv.cfg.bindAddress,
		BaseContext: func(l net.Listener) context.Context {
			return ctx
		},
		Handler: srv.router,
	}

	return nil
}

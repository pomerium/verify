package verify

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

// Server is the verify server backend.
type Server struct {
	cfg *config

	client  *firestore.Client
	http    *http.Server
	router  chi.Router
	storage *Storage
}

// New creates a new Server.
func New(options ...Option) *Server {
	cfg := getConfig(options...)
	return &Server{
		cfg: cfg,
	}
}

// Run runs the server.
func (srv *Server) Run(ctx context.Context) error {
	err := srv.init(ctx)
	if err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		log.Info().
			Str("bind-addr", srv.cfg.bindAddress).
			Msg("starting http server")
		return srv.http.ListenAndServe()
	})
	return eg.Wait()
}

func (srv *Server) init(ctx context.Context) error {
	var err error

	log.Info().
		Str("project-id", srv.cfg.firestoreProjectID).
		Msg("connecting to firestore")
	srv.client, err = firestore.NewClient(ctx, srv.cfg.firestoreProjectID)
	if err != nil {
		return fmt.Errorf("failed to create firestore client: %w", err)
	}
	srv.storage = &Storage{client: srv.client}

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

package verify

import (
	"context"
	"crypto/x509"
	"net"
	"net/http"
	"os"

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

	var verifierOpts tlsVerifierOptions
	if len(cfg.extraCACerts) > 0 {
		pool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load system CA certs")
		}
		for _, certPath := range cfg.extraCACerts {
			cert, err := os.ReadFile(certPath)
			if err != nil {
				log.Fatal().Err(err).Str("path", certPath).Msg("failed to read CA cert")
			} else {
				log.Info().Str("path", certPath).Msg("adding extra CA cert")
			}
			ok := pool.AppendCertsFromPEM(cert)
			if !ok {
				log.Warn().Str("path", certPath).Msg("no CA certs found in file")
			}
		}
		verifierOpts.rootCAs = pool
	}
	return &Server{
		cfg:         cfg,
		tlsVerifier: newTLSVerifier(verifierOpts),
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

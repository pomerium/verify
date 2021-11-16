package verify

import "cloud.google.com/go/firestore"

// config defaults
var (
	DefaultBindAddress  = ":8000"
	DefaultJWKSEndpoint = "" // use the audience
	DefaultProjectID    = firestore.DetectProjectID
)

type config struct {
	bindAddress        string
	firestoreProjectID string
	jwksEndpoint       string
}

// An Option customizes the config.
type Option func(cfg *config)

// WithBindAddress sets the bind address in the config.
func WithBindAddress(bindAddress string) Option {
	return func(cfg *config) {
		cfg.bindAddress = bindAddress
	}
}

// WithJWKSEndpoint sets the jwks endpoint in the config.
func WithJWKSEndpoint(jwksEndpoint string) Option {
	return func(cfg *config) {
		cfg.jwksEndpoint = jwksEndpoint
	}
}

// WithFirestoreProjectID sets the firestore project id in the config.
func WithFirestoreProjectID(projectID string) Option {
	return func(cfg *config) {
		cfg.firestoreProjectID = projectID
	}
}

func getConfig(options ...Option) *config {
	cfg := new(config)
	WithBindAddress(DefaultBindAddress)(cfg)
	// by default the firestore project id is derived from the environment
	WithFirestoreProjectID(DefaultProjectID)(cfg)
	WithJWKSEndpoint(DefaultJWKSEndpoint)(cfg)
	for _, option := range options {
		option(cfg)
	}
	return cfg
}

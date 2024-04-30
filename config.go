package verify

import "cloud.google.com/go/firestore"

// config defaults
var (
	DefaultBindAddress  = ":8000"
	DefaultJWKSEndpoint = "" // use the audience
	DefaultProjectID    = firestore.DetectProjectID
)

type config struct {
	bindAddress         string
	firestoreProjectID  string
	jwksEndpoint        string
	expectedJWTIssuer   string
	expectedJWTAudience string
	extraCACerts        []string
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

// WithExpectedJWTIssuer sets the expected JWT issuer claim in the config. If
// set to the empty string, the issuer claim will not be validated.
func WithExpectedJWTIssuer(issuer string) Option {
	return func(cfg *config) {
		cfg.expectedJWTIssuer = issuer
	}
}

// WithExpectedJWTAudience sets the expected JWT audience claim in the config.
// If set to the empty string, the audience claim will not be validated.
func WithExpectedJWTAudience(audience string) Option {
	return func(cfg *config) {
		cfg.expectedJWTAudience = audience
	}
}

// WithFirestoreProjectID sets the firestore project id in the config.
func WithFirestoreProjectID(projectID string) Option {
	return func(cfg *config) {
		cfg.firestoreProjectID = projectID
	}
}

// WithExtraCACerts adds paths to custom CA certificates to the config.
// Certificates added with this option will be used in addition to the system
// default pool.
func WithExtraCACerts(paths ...string) Option {
	return func(cfg *config) {
		cfg.extraCACerts = append(cfg.extraCACerts, paths...)
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

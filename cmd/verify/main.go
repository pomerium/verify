package main

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/pomerium/verify"
)

func main() {
	addr := verify.DefaultBindAddress
	if v, ok := os.LookupEnv("ADDR"); ok {
		addr = v
	} else if v, ok := os.LookupEnv("PORT"); ok {
		addr = ":" + v
	}

	jwksEndpoint := verify.DefaultJWKSEndpoint
	if v, ok := os.LookupEnv("JWKS_ENDPOINT"); ok {
		jwksEndpoint = v
	}

	firestoreProjectID := verify.DefaultProjectID
	if v, ok := os.LookupEnv("GCLOUD_PROJECT"); ok {
		firestoreProjectID = v
	}

	srv := verify.New(
		verify.WithBindAddress(addr),
		verify.WithFirestoreProjectID(firestoreProjectID),
		verify.WithJWKSEndpoint(jwksEndpoint),
		verify.WithExpectedJWTIssuer(os.Getenv("EXPECTED_JWT_ISSUER")),
		verify.WithExpectedJWTAudience(os.Getenv("EXPECTED_JWT_AUDIENCE")),
	)
	err := srv.Run(context.Background())
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

[![Coverage Status](https://coveralls.io/repos/github/pomerium/verify/badge.svg)](https://coveralls.io/github/pomerium/verify)

# Pomerium Verify service

This example service uses the Pomerium
[Go SDK](https://github.com/pomerium/sdk-go) to parse and display the contents
of the `X-Pomerium-Jwt-Assertion` header. This can help to validate that a
Pomerium deployment is working as expected.

This service is hosted at
[https://verify.pomerium.com](https://verify.pomerium.com), or you can deploy
an instance in your own Pomerium setup.

## Configuration options

The service can be configured with the following environment variables:

- `ADDR`

      Listen address for the service. If neither `ADDR` nor `PORT` is set, the

  service will listen at `:8000`.

- `PORT`

      Listen address port for the service. If neither `ADDR` nor `PORT` is set,

  the service will listen at `:8000`.

- `JWKS_ENDPOINT`

      Allows setting a static URL to use for fetching the public key(s) for

  verifying the Pomerium JWT. If unset, keys will be fetched from the domain
  specified in the JWT `iss` claim (using the internal Pomerium endpoint at
  `/.well-known/pomerium/jwks.json`). Note: in order for this to work correctly,
  you must define `signing_key` or `signing_key_file` in the Pomerium
  configuration.

- `EXPECTED_JWT_ISSUER`

      When set, JWT verification will additionally validate that the issuer claim

  (`iss`) matches the given value.

- `EXPECTED_JWT_AUDIENCE`

      When set, JWT verification will additionally validate that the audience

  claim (`aud`) matches the given value.

- `GCLOUD_PROJECT`

      When set to a Firebase project ID, the service will use [Cloud

  Firestore](https://firebase.google.com/docs/firestore) as a storage backend for
  WebAuthn-related storage. (By default, the service will store this data in
  memory instead.)

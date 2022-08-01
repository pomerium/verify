FROM node:18@sha256:ebd1096a66c724af78abb11e6c81eb05b85fcbe8920af2c24d42b6df6aab2687 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.18-buster@sha256:6960d62610b18b7224d2c5572b4bb177890b9ab7bf70ebaf34e2e9ca662a46e9 as build
WORKDIR /build

COPY Makefile ./Makefile

# download go dependencies
COPY go.mod go.sum ./
RUN go mod download

# build console
COPY --from=ui /build/ui/dist ./ui/dist
COPY ./cmd/ ./cmd/
COPY ./internal/ ./internal/
COPY ./*.go ./
RUN make build-verify

FROM gcr.io/distroless/base-debian10:debug@sha256:d4f8f92882d49b4e0e407da43b7607c6ef3bfb6d8db46a8b9a8cd4064acf4f61
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

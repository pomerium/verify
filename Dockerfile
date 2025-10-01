FROM node:lts-bookworm@sha256:4973262386dc1cb70f7d6fc48a855027d8f12d2d3b1fe559b9db9a4fcb74668f AS ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.25.1-bookworm@sha256:028d35aae29e3e768a7330f4622af402b45e0be58ff5acd49e1eb4f9940f83e0 AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:0d30cde8cd133aa8b8a967bb22bdb3fbfeb4daf413cce46b7149c4d4672dfe1e
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

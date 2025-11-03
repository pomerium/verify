FROM node:lts-bookworm@sha256:34af25027ee1b8bffd482ba995ec1e577fbd398db87beb4c60b80c2c9c025127 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.25.3-bookworm@sha256:ee420c17fa013f71eca6b35c3547b854c838d4f26056a34eb6171bba5bf8ece4 AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:d4bcaaac2088ef3bb6ca9f600cfd3f34939fb8fb4658243ee1b00c309f509eb7
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

FROM node:lts-bookworm@sha256:5f21943fe97b24ae1740da6d7b9c56ac43fe3495acb47c1b232b0a352b02a25c as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.21.5-bookworm@sha256:a6b787c7f9046e3fdaa97bca1f76fd23ff4108f612de885e1af87e0dccc02f99 as build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:5e24c7a60ad746d78fd96034b6d043c00ef6ed94ec55ee7882d93162c939f3a1
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

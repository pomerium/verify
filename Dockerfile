FROM node:lts-bookworm@sha256:f3299f16246c71ab8b304d6745bb4059fa9283e8d025972e28436a9f9b36ed24 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.22.0-bookworm@sha256:925fe3fa28ba428cf67a7947ae838f8a1523117b40e3e6b5106c378e3f97fa29 as build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:1d91d5f201f4be1ab11bdbcd3ef9d1a42bdfd29863bb007eba11eaa55172f85d
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

FROM node:22.22.2-bookworm@sha256:ecabd1cb6956d7acfffe8af6bbfbe2df42362269fd28c227f36367213d0bb777 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.25.9-bookworm@sha256:29e59af995c51a5bf63d072eca973b918e0e7af4db0e4667aa73f1b8da1a6d8c AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:1f8759794cab46f0673e14afc03e3623cbd803b683abf7e3143fd041cc2e89f7
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

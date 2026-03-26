FROM node:22.22.2-bookworm@sha256:7e791fc54bd02fc89fd4fb39eb37e5bea753c75679c8022478d81679367d995a AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.25.8-bookworm@sha256:7fb09d8804035fbde8a84ed59ca9f46dd68c6f160f9d193e98d795d8d9e002ec AS build
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

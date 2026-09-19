FROM node:26.8.2-bookworm@sha256:e7bc1a4cd2419953c91f9a6f7bb6efb3737773093fb4ded0b1c77a0a5831fac4 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.27.1-bookworm@ AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:0a5b69081973d7069cbeab49c24851dd517887767e934fc8e6f808ea08770ef4
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

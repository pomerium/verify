FROM node:26.13.1-bookworm@sha256:893ea48106e988a5f91e69629d6472600daa9d0494adf7db87eb0b886ce82255 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.27.4-bookworm@sha256:ce59b41415e2524ae99926fe1f37ce48ba8c194dcbd2a56a65436cdcefa2bd34 AS build
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

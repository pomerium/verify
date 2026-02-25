FROM node:22.22.0-bookworm@sha256:a871fb3fb50960e4701335cf5aa3ee7a1c6f966127ddc5d9b9a6035d58f9450f AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.25.7-bookworm@sha256:ef2563a2e7b73a72b643c914a60189ab8273080c715506326434873ab7d6cce8 AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:e8075f7da06319e4ac863d31fa11354003c809ef9f1b52fe32ef39e876ac16c5
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

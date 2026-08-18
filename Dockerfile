FROM node:26.7.0-bookworm@sha256:0353e48e0e8a993db87b720c242f54b207059d1bcc0106534896e8a11054c837 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.26.6-bookworm@sha256:116d58cbd88c1297624acc6e967a060012422bacf9930927e23fb719189c6f36 AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:db9a19954a9f79ef814bb204aa7788cfcfd3e00a4bfd755ad93d2e88721125b6
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

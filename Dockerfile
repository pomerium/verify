FROM node:22.22.1-bookworm@sha256:b501c082306a4f528bc4038cbf2fbb58095d583d0419a259b2114b5ac53d12e9 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.25.8-bookworm@sha256:45aaee97bdc61507976d86b12f65fa8394e107caba1d13b961a15f63edf24ecd AS build
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

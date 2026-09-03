FROM node:26.8.1-bookworm@sha256:9f94d34c787165dca03b74e5bf9c3bf90e8de79b19aa3d87fe1fa1694bf75c89 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.27.1-bookworm@sha256:648f440f42a0958804efb24df176f806f9d353b41f1c0627f666428e40310f6b AS build
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

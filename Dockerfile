FROM node:19@sha256:d903b23148dccca63152cb2bdf6fdfa062842dcce49f0d45558ca06a8ab14598 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.20.1-buster@sha256:d99d3614cfc6bc0f7e4d96b54682d403f2a7c2ab82c6001a7701e937dd8f4d6d as build
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

FROM gcr.io/distroless/base-debian10:debug@sha256:0828a5cfb40820d9110d4fb379edfca129ff5ec01d32d4f8acfd60173e2c49f3
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

FROM node:lts-bookworm@sha256:fd0115473b293460df5b217ea73ff216928f2b0bb7650c5e7aa56aae4c028426 as ui
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

FROM node:18@sha256:fb040635dd58488aaa66d75e253997e414e1c4adf852ce92c6a8bb2f35e8e24b as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.18.4-buster@sha256:477b10a289d7be9c5124da52f2ff86d13c24e591b82da3b6b1ad3b8f1a04c8d9 as build
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

FROM gcr.io/distroless/base-debian10:debug@sha256:d4f8f92882d49b4e0e407da43b7607c6ef3bfb6d8db46a8b9a8cd4064acf4f61
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

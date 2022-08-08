FROM node:18@sha256:a6f295c2354992f827693a2603c8b9b5b487db4da0714f5913a917ed588d6d41 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.19.0-buster@sha256:a7a23f1fba8390b1e038f017c85259c878f406301643653ec6e5b97e75668789 as build
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

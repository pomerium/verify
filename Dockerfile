FROM node:19@sha256:f62cd4433edc4b397218d9bc005a2f2865d7827753fd480a166e247127384900 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.19.5-buster@sha256:150a98859da4c88935958d21d4fca1c74a139eac88e1eb2364c8f60d69236ede as build
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

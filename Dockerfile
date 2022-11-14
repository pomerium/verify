FROM node:19@sha256:6da4e30e3952e460fe4ad256e46a8b79acce46dd596bbe4ef882d5ec0d1ef6cb as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.19.3-buster@sha256:0092b239c5b5dc572729d8fd3c1842513ed1222502f86be65cbb875c5b9b1e8a as build
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

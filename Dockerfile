FROM node:18@sha256:8a45c95c328809e7e10e8c9ed5bf8374620d62e52de1df7ef8e71a9596ec8676 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.19.0-buster@sha256:d84495e2ad12dfeb51dec3c9f170659ebd09db234d6990b3364f877785684a14 as build
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

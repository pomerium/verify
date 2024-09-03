FROM node:lts-bookworm@sha256:a4d1de4c7339eabcf78a90137dfd551b798829e3ef3e399e0036ac454afa1291 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.23.0-bookworm@sha256:31dc846dd1bcca84d2fa231bcd16c09ff271bcc1a5ae2c48ff10f13b039688f3 as build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:af772ed0ce52d8994acedc3ec84a9d22e9366dda8767f17d1bb2213b06beaff5
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

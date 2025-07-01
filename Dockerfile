FROM node:lts-bookworm@sha256:0c0734eb7051babbb3e95cd74e684f940552b31472152edf0bb23e54ab44a0d7 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.24.4-bookworm@sha256:940ac576af6f5d674dd5a173ee0d93cd8ad317890e823ae309f9c2c7b8fa788c AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:7d1d72086ccf7b5c7e0f612dd59ae064765a529daafaecac97ea4a8b48b69e93
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

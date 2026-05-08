FROM node:22.22.2-bookworm@sha256:9059d9d7db987b86299e052ff6630cd95e5a770336967c21110e53289a877433 AS ui
WORKDIR /build

COPY Makefile ./Makefile

# build ui
COPY ./ui/ ./ui/
RUN make npm-install
RUN make build-ui

FROM golang:1.26.3-bookworm@sha256:b09e568dcf2a1ff3ce09a230ea234193fc014dc195472fe63316e50238453d96 AS build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:fd8e2df55ce1400d6f651f14a030fd38868fca0e1ab6989bc1641cd5fc0f3335
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

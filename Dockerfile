FROM node:lts-bookworm@sha256:db5dd2f30cb82a8bdbd16acd4a8f3f2789f5b24f6ce43f98aa041be848c82e45 as ui
WORKDIR /build

COPY Makefile ./Makefile

# download yarn dependencies
COPY ui/yarn.lock ./ui/yarn.lock
COPY ui/package.json ./ui/package.json
RUN make yarn

# build ui
COPY ./ui/ ./ui/
RUN make build-ui

FROM golang:1.23.1-bookworm@sha256:dba79eb312528369dea87532a65dbe9d4efb26439a0feacc9e7ac9b0f1c7f607 as build
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

FROM gcr.io/distroless/base-debian12:debug@sha256:662eaa2606087124ac1fc108291ecad341f6376ce6fa28ac7e1655ec76c6e6d9
COPY --from=build /build/bin/* /bin/
ENTRYPOINT ["/bin/verify"]

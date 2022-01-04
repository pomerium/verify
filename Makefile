NAME := verify
PKG := github.com/pomerium/$(NAME)
PREFIX?=$(shell pwd)
BUILDDIR := ${PREFIX}/dist
BINDIR := ${PREFIX}/bin

GO ?= "go"
YARN ?= "yarn"

# Build binary artifact
.PHONY: build
build: build-ui build-verify

# Build verify executable
.PHONY: build-verify
build-verify:
	@echo "==> $@"
	$(GO) build -o ${BINDIR}/${NAME} cmd/verify/*.go

# Build frontend javascript
.PHONY: build-ui
build-ui: yarn
	@echo "==> $@"
	cd ui; $(YARN) build

# Initial yarn install
.PHONY: yarn
yarn:
	@echo "==> $@"
	cd ui; $(YARN)

# Run go tests
test:
	@echo "==> $@"
	go test -v ./...

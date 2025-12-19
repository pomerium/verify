NAME := verify
PKG := github.com/pomerium/$(NAME)
PREFIX?=$(shell pwd)
BUILDDIR := ${PREFIX}/dist
BINDIR := ${PREFIX}/bin

.PHONY: all
all: build lint test

# Build binary artifact
.PHONY: build
build: build-ui build-verify

# Build verify executable
.PHONY: build-verify
build-verify:
	@echo "==> $@"
	go build -o ${BINDIR}/${NAME} cmd/verify/*.go

# Build frontend javascript
.PHONY: build-ui
build-ui: npm-install
	@echo "==> $@"
	cd ui; npm run build

# Initial yarn install
.PHONY: npm-install
npm-install:
	@echo "==> $@"
	cd ui; npm ci

# Run go tests
.PHONY: test
test:
	@echo "==> $@"
	go test -v ./...

.PHONY: cover
cover: ## Runs go test with coverage
	@echo "==> $@"
	@go test -race -coverprofile=coverage.txt ./...

.PHONY: lint
lint:
	@echo "@==> $@"
	golangci-lint run --fix ./...

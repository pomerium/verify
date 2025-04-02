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
	cd ui; $(YARN) --network-timeout 100000

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
	@VERSION=$$(go run github.com/mikefarah/yq/v4@v4.34.1 '.jobs.lint.steps[] | select(.uses == "golangci/golangci-lint-action*") | .with.version' .github/workflows/lint.yml) && \
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@$$VERSION run ./...

.DEFAULT_GOAL := help
SHELL := /bin/bash

APP_NAME := stylist
BUILD_DIR := ./dist
DST_DIR := /usr/local/bin

##@ App

.PHONY: coverage
coverage: gocovsh ## Show code coverage
	@make test
	gocovsh --profile coverage.out

.PHONY: build
build: ## Build the app
	go mod tidy
	go build -o ${BUILD_DIR}/${APP_NAME} .

.PHONY: install
install: build ## Install the app
	install -d ${DST_DIR}
	install -m755 ${BUILD_DIR}/${APP_NAME} ${DST_DIR}/

.PHONY: generate
generate: go-enum
	go generate ./...

.PHONY: lint
lint: actionlint golangci-lint ## Lint the app
	actionlint
	golangci-lint run

.PHONY: format
format: pin-github-action ## Format the app
	gofmt -w .
	pin-github-action .github/workflows/*.yml

.PHONY: test
test: ## Test the app
	go mod tidy
	go test --coverprofile=coverage.out ./...


##@ Other

.PHONY: setup
setup: ## Bootstrap for local development
	@if ! command -v gh >/dev/null 2>&1; then echo "Unable to find gh!"; exit 1; fi
	@if ! command -v git >/dev/null 2>&1; then echo "Unable to find git!"; exit 1; fi
	@if ! command -v go >/dev/null 2>&1; then echo "Unable to find go!"; exit 1; fi

# Via https://www.thapaliya.com/en/writings/well-documented-makefiles/
# Note: The `##@` comments determine grouping
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
	@echo ""


# Dependencies

.PHONY: actionlint
actionlint:
	@if ! command -v actionlint >/dev/null 2>&1; then go install github.com/rhysd/actionlint/cmd/actionlint@latest; fi

.PHONY: go-enum
go-enum:
	@if ! command -v go-enum >/dev/null 2>&1; then go install github.com/abice/go-enum@latest; fi

.PHONY: gocovsh
gocovsh:
	@if ! command -v gocovsh >/dev/null 2>&1; then go install github.com/orlangure/gocovsh@latest; fi

.PHONY: golangci-lint
golangci-lint:
	@if ! command -v golangci-lint >/dev/null 2>&1; then go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; fi

.PHONY: pin-github-action
pin-github-action:
	@if ! command -v pin-github-action >/dev/null 2>&1; then npm install -g pin-github-action; fi

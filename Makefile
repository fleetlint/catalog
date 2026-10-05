# Task-runner contract. `make check` is the definition of green.
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL:
.DEFAULT_GOAL := check-fast

GOLANGCI_VERSION    ?= v2.14.0
GOVULNCHECK_VERSION ?= v1.8.0
# Tools are built with this module's toolchain: golangci-lint refuses code that targets a newer Go than it was built with.
TOOLCHAIN := $(shell go env GOVERSION)

.PHONY: help tools fmt lint test cover audit build check-fast check

help: ## list targets
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | awk -F':.*##' '{printf "  %-12s %s\n", $$1, $$2}'

tools: ## install the pinned tools (gitleaks and fleetlint come from their own installers)
	GOTOOLCHAIN=$(TOOLCHAIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	GOTOOLCHAIN=$(TOOLCHAIN) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

fmt: ## format in place
	golangci-lint fmt

lint: ## format check, linters, vet
	golangci-lint fmt --diff
	golangci-lint run
	go vet ./...

test: ## the presence test for presets and templates
	go test -race ./...

cover: test ## nothing to measure: the module is data and one embed declaration

audit: ## vulnerabilities, tidy module graph, secrets
	govulncheck ./...
	go mod tidy -diff
	gitleaks git --no-banner --redact .

build: ## compile the module
	go build ./...

check-fast: lint test ## before every commit

check: lint test cover audit build ## definition of green
	fleetlint check --fail-on error

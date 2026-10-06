# Task-runner contract. `make check` is the definition of green.
SHELL := bash
.SHELLFLAGS := -euo pipefail -c
.ONESHELL:
.DEFAULT_GOAL := check-fast

# renovate: datasource=go depName=github.com/golangci/golangci-lint/v2
GOLANGCI_VERSION    ?= v2.14.0
# renovate: datasource=go depName=golang.org/x/vuln
GOVULNCHECK_VERSION ?= v1.8.0
# Tools are built with this module's toolchain: golangci-lint refuses code that targets a newer Go than it was built with.
TOOLCHAIN := $(shell go env GOVERSION)

.PHONY: help tools fmt lint test cover audit build check-fast check gen

help: ## list targets
	@grep -E '^[a-z-]+:.*##' $(MAKEFILE_LIST) | awk -F':.*##' '{printf "  %-12s %s\n", $$1, $$2}'

tools: ## install the pinned tools (gitleaks and fleetlint come from their own installers)
	GOTOOLCHAIN=$(TOOLCHAIN) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_VERSION)
	GOTOOLCHAIN=$(TOOLCHAIN) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)

fmt: ## format in place
	golangci-lint fmt

lint: ## format check, linters, vet, strict preset current
	golangci-lint fmt --diff
	golangci-lint run
	go vet ./...
	@cp presets/strict.yaml /tmp/strict.before && python3 scripts/gen-strict.py >/dev/null && \
	  (diff -q /tmp/strict.before presets/strict.yaml >/dev/null || { echo "presets/strict.yaml was stale; regenerated, review and commit"; exit 1; })

gen: ## regenerate presets/strict.yaml from minimal and recommended
	python3 scripts/gen-strict.py

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

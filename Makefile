# ==============================================================================
# Noticeal - Development Makefile
# ==============================================================================

.DEFAULT_GOAL := help

GREEN  := \033[0;32m
YELLOW := \033[1;33m
BLUE   := \033[0;34m
RED    := \033[0;31m
RESET  := \033[0m

.PHONY: help run fmt vet test build install clean \
        release-build release-snapshot release-check \
        doctor version tidy update

help: ## Show available commands
	@echo ""
	@echo "$(BLUE)Noticeal Development Commands$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[32m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# ==============================================================================
# Development
# ==============================================================================

run: ## Run Noticeal
	go -C app run ./cmd

build: ## Build local binary
	@mkdir -p bin
	go -C app build -o ../bin/noticeal ./cmd
	@echo "$(GREEN)✓ Binary generated in bin/noticeal$(RESET)"

install: build ## Install Noticeal locally
	sudo install -m 755 bin/noticeal /usr/local/bin/noticeal
	@echo "$(GREEN)✓ Installed in /usr/local/bin/noticeal$(RESET)"

version: ## Display installed version
	noticeal --version

# ==============================================================================
# Quality
# ==============================================================================

fmt: ## Format source code
	go -C app fmt ./...

vet: ## Run go vet
	go -C app vet ./...

test: ## Run unit tests
	go -C app test ./...

check: fmt vet test ## Run all quality checks

# ==============================================================================
# Dependencies
# ==============================================================================

tidy: ## Clean go.mod / go.sum
	go -C app mod tidy

update: ## Update dependencies
	go -C app get -u ./...
	go -C app mod tidy

# ==============================================================================
# Release
# ==============================================================================

release-build: ## Build release artifacts (no GitHub release)
	goreleaser build --clean

release-snapshot: ## Simulate a release locally
	goreleaser release --snapshot --clean

release-check: ## Validate GoReleaser configuration
	goreleaser check

# ==============================================================================
# Utilities
# ==============================================================================

clean: ## Remove generated files
	rm -rf bin
	rm -rf dist

doctor: ## Display development environment
	@echo ""
	@echo "$(BLUE)Environment$(RESET)"
	@echo ""
	@go version
	@git --version
	@goreleaser --version
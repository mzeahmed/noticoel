# ==============================================================================
# Noticoel - Development Makefile
# ==============================================================================

.DEFAULT_GOAL := help

GREEN  := \033[0;32m
YELLOW := \033[1;33m
BLUE   := \033[0;34m
RED    := \033[0;31m
RESET  := \033[0m

# ==============================================================================
# Environment
# ==============================================================================

-include .env
export

.PHONY: help \
        run dev build install \
        fmt vet lint test check \
        tidy update \
        migrate-up migrate-down migrate-status migrate-create sqlc \
        send \
        release-build release-snapshot release-check \
        doctor version clean

help: ## Show available commands
	@echo ""
	@echo "$(BLUE)Noticoel Development Commands$(RESET)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "  \033[32m%-22s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""

# ==============================================================================
# Development
# ==============================================================================

run: ## Run Noticoel
	go -C app run ./cmd

dev: ## Run noticoel with hot reload (Air)
	air

build: ## Build local binary
	@mkdir -p bin
	go -C app build -o ../bin/noticoel ./cmd
	@echo "$(GREEN)✓ Binary generated in bin/noticoel$(RESET)"

#install: build ## Install noticoel locally
	#sudo install -m 755 bin/noticoel /usr/local/bin/noticoel
	@#echo "$(GREEN)✓ Installed in /usr/local/bin/noticoel$(RESET)"

# ==============================================================================
# Quality
# ==============================================================================

fmt: ## Format source code
	go -C app fmt ./...

vet: ## Run go vet
	go -C app vet ./...

lint: ## Run golangci-lint
	golangci-lint run

test: ## Run unit tests
	go -C app test ./...

check: fmt vet lint test ## Run all quality checks

# ==============================================================================
# Dependencies
# ==============================================================================

tidy: ## Clean go.mod / go.sum
	go -C app mod tidy

update: ## Update dependencies
	go -C app get -u ./...
	go -C app mod tidy

# ==============================================================================
# Database
# ==============================================================================

migrate-up: ## Apply all database migrations
	goose -dir app/migrations sqlite3 data/noticoel.db up

migrate-down: ## Roll back the last migration
	goose -dir app/migrations sqlite3 data/noticoel.db down

migrate-status: ## Show migration status
	goose -dir app/migrations sqlite3 data/noticoel.db status

migrate-create: ## Create a new SQL migration (NAME=create_events_table)
	@test -n "$(NAME)" || (echo "Usage: make migrate-create NAME=create_events_table" && exit 1)
	goose -dir app/migrations create $(NAME) sql

sqlc: ## Regenerate Go code from SQL queries
	cd app && sqlc generate

# ==============================================================================
# Examples
# ==============================================================================

send: ## Send a sample event (EVENT=workflow-success|workflow-failure|release, default workflow-success)
	bash examples/scripts/send.sh examples/events/$(or $(EVENT),workflow-success).json

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

doctor: ## Display development environment
	@echo ""
	@echo "$(BLUE)Development Environment$(RESET)"
	@echo ""

	@printf "%-18s" "Go:"
	@go version

	@printf "%-18s" "Git:"
	@git --version

	@printf "%-18s" "Air:"
	@air -v || echo "Not installed"

	@printf "%-18s" "Goose:"
	@goose -version || echo "Not installed"

	@printf "%-18s" "golangci-lint:"
	@golangci-lint version || echo "Not installed"

	@printf "%-18s" "GoReleaser:"
	@goreleaser --version

version: ## Show installed tool versions
	@go version
	@air -v || true
	@goose -version || true
	@golangci-lint version || true
	@goreleaser --version

clean: ## Remove generated files
	rm -rf bin
	rm -rf dist
	rm -rf tmp
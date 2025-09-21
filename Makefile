# Simple Makefile for a Go project

# Tools
GINKGO := $(shell go env GOPATH)/bin/ginkgo
MIGRATIONS_DIR := ./src/migrations

# --- Env loading ---
# Load environment variables from .env automatically for CLI targets
ENV_FILE := ./src/.env
LOAD_ENV := set -a; if [ -f $(ENV_FILE) ]; then . $(ENV_FILE); fi; set +a

# Build the application
all: build-binary test

build:
	@echo "Building binary for air..."
	@cd src && go build -o main ./cmd/api

build-binary:
	@echo "Building binary..."
	@mkdir -p bin
	@go build -o bin/api ./src/cmd/api

# Run the application
run:
	@go run ./src/cmd/api

# --- Migrations (Goose) ---
# Generate new Go migration skeleton (requires name=...)
makemigrations:
	@if [ -z "$(name)" ]; then \
		echo "Usage: make makemigrations name=<migration_name>"; \
		exit 1; \
	fi
	@echo "Creating migration: $(name)"
	@go run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATIONS_DIR) create $(name) go

# Run migrations using our custom binary (imports Go migrations)
migrate:
	@$(LOAD_ENV); \
	if [ -z "$$DB_USERNAME" ] || [ -z "$$DB_PASSWORD" ] || [ -z "$$DB_HOST" ] || [ -z "$$DB_PORT" ] || [ -z "$$DB_DATABASE" ] || [ -z "$$DB_SCHEMA" ]; then \
		echo "Set DB_* env vars before running migrate"; \
		exit 1; \
	fi
	@go run ./src/cmd/migrate up

migrate-down:
	@$(LOAD_ENV); \
	if [ -z "$$DB_USERNAME" ] || [ -z "$$DB_PASSWORD" ] || [ -z "$$DB_HOST" ] || [ -z "$$DB_PORT" ] || [ -z "$$DB_DATABASE" ] || [ -z "$$DB_SCHEMA" ]; then \
		echo "Set DB_* env vars before running migrate-down"; \
		exit 1; \
	fi
	@go run ./src/cmd/migrate down

migrate-status:
	@$(LOAD_ENV); \
	if [ -z "$$DB_USERNAME" ] || [ -z "$$DB_PASSWORD" ] || [ -z "$$DB_HOST" ] || [ -z "$$DB_PORT" ] || [ -z "$$DB_DATABASE" ] || [ -z "$$DB_SCHEMA" ]; then \
		echo "Set DB_* env vars before running migrate-status"; \
		exit 1; \
	fi
	@go run ./src/cmd/migrate status

dev-up:
	@docker compose -f src/docker-compose.yml up -d

dev-rebuild:
	@docker compose -f src/docker-compose.yml up -d --build

dev-down:
	@docker compose -f src/docker-compose.yml down 

dev-logs:
	@docker compose -f src/docker-compose.yml logs -f api

# Test the application
test:
	@echo "Testing (fast specs)..."
	@cd src && $(GINKGO) -r -p -v --skip-package=internal/database

test-db:
	@echo "Testing (db)..."
	@cd src && $(GINKGO) -v ./internal/database

test-all:
	@echo "Testing (all specs)..."
	@make test
	@make test-db

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f bin/api

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            cd src && air; \
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                cd src && air; \
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build-binary run test clean watch dev-up dev-down dev-rebuild dev-logs test-db test-all makemigrations migrate migrate-down migrate-status
.PHONY: help build run test test-coverage clean migrate-up migrate-down sqlc docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building..."
	@go build -o bin/chirpy ./cmd/api

run: ## Run the application
	@echo "Running..."
	@go run ./cmd/api/main.go

dev: ## Run with hot reload (requires air)
	@air

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html

sqlc: ## Generate sqlc code from SQL queries
	@echo "🔄 Generating sqlc code from SQL queries..."
	@sqlc generate
	@echo "✅ Generated code in internal/database/"
	@echo "⚠️  DO NOT manually edit files in internal/database/"
	@echo "💡 To add queries: edit sql/queries/*.sql then run 'make sqlc'"

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@goose -dir sql/schema postgres "$$DB_URL" up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@goose -dir sql/schema postgres "$$DB_URL" down

migrate-status: ## Check migration status
	@goose -dir sql/schema postgres "$$DB_URL" status

docker-up: ## Start Docker containers
	@docker-compose up -d

docker-down: ## Stop Docker containers
	@docker-compose down

lint: ## Run linter
	@golangci-lint run

format: ## Format code
	@go fmt ./...

tidy: ## Tidy go modules
	@go mod tidy

install-tools: ## Install development tools
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
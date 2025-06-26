# Pet Store Backend Makefile

.PHONY: help build test test-unit test-integration test-coverage clean docker-build docker-up docker-down lint fmt deps check

# Default target
help: ## Show this help message
	@echo 'Usage:'
	@echo '  make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Build targets
build: ## Build the application
	@echo "Building application..."
	@go build -o bin/server ./cmd/server

# Test targets
test: test-unit ## Run all tests
	@echo "All tests completed"

test-unit: ## Run unit tests
	@echo "Running unit tests..."
	@go test -v -race -short ./...

test-integration: ## Run integration tests (requires test database)
	@echo "Running integration tests..."
	@go test -v -race -tags=integration ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w .

# Dependency management
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Validation targets
check: fmt lint test-unit ## Run all code quality checks

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker-compose build backend

docker-up: ## Start services with Docker Compose
	@echo "Starting services..."
	@docker-compose up -d

docker-down: ## Stop services
	@echo "Stopping services..."
	@docker-compose down

docker-logs: ## Show service logs
	@docker-compose logs -f backend

# Database targets
db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	@go run ./cmd/server migrate

db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "Resetting database..."
	@docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS petstore; CREATE DATABASE petstore;"
	@$(MAKE) db-migrate

# Generate targets
generate: ## Generate GraphQL code
	@echo "Generating GraphQL code..."
	@go generate ./...

# Benchmark targets
bench: ## Run benchmarks
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Clean targets
clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@docker-compose down --volumes --remove-orphans

# Development targets
dev: docker-up ## Start development environment
	@echo "Development environment started"
	@echo "GraphQL Playground: http://localhost:8080/playground"
	@echo "API Endpoint: http://localhost:8080/graphql"

# Security targets
security: ## Run security scan
	@echo "Running security scan..."
	@gosec ./...

# Documentation targets
docs: ## Generate documentation
	@echo "Generating documentation..."
	@godoc -http=:6060 &
	@echo "Documentation available at: http://localhost:6060"

# Install development tools
install-tools: ## Install development tools
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@go install golang.org/x/tools/cmd/godoc@latest
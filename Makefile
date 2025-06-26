# Pet Store Project Makefile
# This Makefile provides commands for the entire pet store project

.PHONY: help setup build test clean docker-build docker-up docker-down generate lint fmt deps dev

# Default target
help: ## Show this help message
	@echo 'Pet Store Project - Available Commands:'
	@echo ''
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development Setup
setup: ## Initial project setup
	@echo "🚀 Setting up Pet Store project..."
	@echo "📦 Installing Go dependencies..."
	cd backend && go mod download
	@echo "🔧 Installing development tools..."
	cd backend && go install github.com/99designs/gqlgen@latest
	cd backend && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "📁 Creating directories..."
	mkdir -p backend/uploads
	mkdir -p backend/bin
	@echo "⚙️  Setting up environment..."
	@if [ ! -f backend/.env ]; then \
		cp backend/.env.example backend/.env; \
		echo "✅ Created backend/.env from example"; \
	else \
		echo "ℹ️  backend/.env already exists"; \
	fi
	@echo ""
	@echo "✅ Setup complete!"
	@echo "📝 Next steps:"
	@echo "   1. Edit backend/.env with your configuration"
	@echo "   2. Run 'make dev' to start development environment"
	@echo "   3. Visit http://localhost:8080/playground for GraphQL playground"

# Code Generation
generate: ## Generate GraphQL code
	@echo "🔄 Generating GraphQL code..."
	cd backend && gqlgen generate
	@echo "✅ Code generation complete"

# Build
build: ## Build backend binary
	@echo "🔨 Building backend..."
	cd backend && go build -o bin/server ./cmd/server
	@echo "✅ Build complete: backend/bin/server"

# Testing
test: ## Run all tests
	@echo "🧪 Running tests..."
	cd backend && go test -v -race ./...
	@echo "✅ Tests complete"

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	cd backend && go test -v -race -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: backend/coverage.html"

# Code Quality
lint: ## Run linter
	@echo "🔍 Running linter..."
	cd backend && GOTOOLCHAIN=local golangci-lint run ./...
	@echo "✅ Linting complete"

fmt: ## Format code
	@echo "✨ Formatting code..."
	cd backend && go fmt ./...
	@echo "✅ Code formatted"

check: ## Run all code quality checks
	@$(MAKE) fmt
	@$(MAKE) lint
	@$(MAKE) test

# Dependencies
deps: ## Update dependencies
	@echo "📦 Updating dependencies..."
	cd backend && go mod download && go mod tidy
	@echo "✅ Dependencies updated"

# Docker Development
docker-build: ## Build Docker images
	@echo "🐳 Building Docker images..."
	docker-compose build
	@echo "✅ Docker images built"

docker-up: ## Start development environment with Docker
	@echo "🚀 Starting development environment..."
	docker-compose up -d
	@echo ""
	@echo "✅ Development environment started!"
	@echo "🌐 Backend API: http://localhost:8080/graphql"
	@echo "🎮 GraphQL Playground: http://localhost:8080/playground"
	@echo "📊 Logs: make docker-logs"

docker-down: ## Stop Docker services
	@echo "🛑 Stopping Docker services..."
	docker-compose down
	@echo "✅ Services stopped"

docker-logs: ## Show Docker logs
	@echo "📋 Showing Docker logs..."
	docker-compose logs -f backend

# Development Shortcuts
dev: docker-up ## Start full development environment
	@echo ""
	@echo "🎉 Development environment ready!"
	@echo ""
	@echo "📝 Quick Links:"
	@echo "   Backend API:        http://localhost:8080/graphql"
	@echo "   GraphQL Playground: http://localhost:8080/playground"
	@echo ""
	@echo "🔧 Useful commands:"
	@echo "   make docker-logs    # View logs"
	@echo "   make docker-down    # Stop services"
	@echo "   make test           # Run tests"

stop: docker-down ## Stop development environment

restart: docker-down docker-up ## Restart development environment

# Database
db-reset: ## Reset database (WARNING: Deletes all data)
	@echo "⚠️  Resetting database (this will delete all data)..."
	@read -p "Are you sure? (y/N): " confirm && [ "$$confirm" = "y" ]
	docker-compose exec postgres psql -U postgres -c "DROP DATABASE IF EXISTS petstore; CREATE DATABASE petstore;"
	@echo "✅ Database reset complete"

# Cleaning
clean: ## Clean all build artifacts and containers
	@echo "🧹 Cleaning up..."
	rm -rf backend/bin/
	rm -rf backend/coverage.out backend/coverage.html
	rm -rf backend/uploads/*
	docker-compose down --volumes --remove-orphans
	docker system prune -f
	@echo "✅ Cleanup complete"

# Security
security-scan: ## Run security scan
	@echo "🔒 Running security scan..."
	cd backend && gosec ./...
	@echo "✅ Security scan complete"


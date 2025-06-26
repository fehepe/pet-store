# 🐾 Pet Store

A complete pet store management system with GraphQL API backend built with Go, PostgreSQL, and Redis.

## 📁 Project Structure

```
pet-store/
├── backend/           # GraphQL API service
│   ├── cmd/          # Application entrypoints
│   ├── internal/     # Private application code
│   │   ├── graph/    # GraphQL schema and resolvers
│   │   ├── service/  # Business logic layer
│   │   ├── repository/ # Data access layer
│   │   ├── models/   # Domain models
│   │   └── auth/     # Authentication middleware
│   └── pkg/          # Public libraries
├── docker-compose.yml # Multi-service orchestration
└── Makefile          # Project automation
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.21+ (for local development)

### Start the Application
```bash
# Clone and start all services
git clone https://github.com/fehepe/pet-store
cd pet-store
docker-compose up -d

# Verify services are running
docker-compose ps
```

### Access Points
- **GraphQL Playground**: http://localhost:8080/playground
- **GraphQL API**: http://localhost:8080/graphql
- **Health Check**: http://localhost:8080/health

## 🏗️ Services

### Backend Service
- **Technology**: Go with GraphQL
- **Features**: 
  - Pet inventory management
  - Store management
  - Customer orders
  - Role-based authentication (Merchant/Customer)
  - Data caching with Redis
  - Encrypted sensitive data
- **Port**: 8080

### Database Services
- **PostgreSQL**: Primary data storage (Port 5432)
- **Redis**: Caching and session management (Port 6379)

## 📚 Documentation

For detailed API documentation, authentication, and usage examples, see the [backend documentation](./backend/README.md).

## 🛠️ Development

```bash
# Build and test
make build
make test

# Development with hot reload
make dev

# Stop all services
docker-compose down
```

## 🔧 Configuration

Environment variables are configured in `docker-compose.yml` and can be customized in `backend/.env` for local development.
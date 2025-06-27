# 🐾 Pet Store

A complete pet store management system with GraphQL API backend built with Go, PostgreSQL, Redis, and a React TypeScript frontend.

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
├── frontend/         # React TypeScript application
│   ├── src/
│   │   ├── components/  # React components
│   │   ├── contexts/    # Context providers (Auth, Cart)
│   │   ├── graphql/     # GraphQL queries and mutations
│   │   └── types/       # TypeScript type definitions
│   └── public/
├── docker-compose.yml # Multi-service orchestration
└── Makefile          # Project automation
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for local development)

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
- **Frontend Application**: http://localhost:3000
- **GraphQL Playground**: http://localhost:8080/playground
- **GraphQL API**: http://localhost:8080/graphql
- **Health Check**: http://localhost:8080/health

## 🏗️ Services

### Frontend Service
- **Technology**: React with TypeScript
- **Features**:
  - Customer authentication
  - Browse available pets by store
  - Instant purchase functionality
  - Shopping cart management
  - Bulk checkout
  - Real-time error handling
  - Responsive design with Material-UI
- **Port**: 3000

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

## 🎯 Usage

### Customer Flow
1. Navigate to http://localhost:3000
2. Login with any username/password and a valid Store ID
3. Browse available pets
4. Either:
   - Click "Buy Now" for instant purchase
   - Add pets to cart and checkout multiple pets at once
5. View purchase confirmations

### Merchant Flow (API Only)
Use the GraphQL Playground at http://localhost:8080/playground with merchant credentials to:
- Create and manage stores
- Add/remove pets from inventory
- View sold and unsold pets
- Track sales by date range

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

# 🐾 Pet Store Application

A complete pet store management system featuring a GraphQL API backend built with Go and a modern React TypeScript frontend.

## 🌟 Features

### ✨ Frontend Features
- **Store Selection**: Dropdown populated from API instead of manual UUID entry
- **Pet Browsing**: View available pets with professional Material-UI design
- **Default Images**: Species-specific default images for pets without pictures
- **Shopping Cart**: Add pets to cart with validation to prevent duplicates
- **Instant Purchase**: Buy pets immediately or checkout multiple pets from cart
- **Real-time Updates**: Live inventory updates and error handling
- **Responsive Design**: Works on desktop, tablet, and mobile devices

### 🔧 Backend Features  
- **GraphQL API**: Modern API with schema-first development
- **Store Management**: Create and manage pet stores
- **Pet Inventory**: Full CRUD operations for pet management
- **Order Processing**: Handle individual and bulk pet purchases
- **Public Endpoints**: Store listing and pet browsing without authentication
- **Data Caching**: Redis integration for improved performance
- **Database Migrations**: Automated schema management with seed data

## 📁 Project Structure

```
pet-store/
├── backend/              # Go GraphQL API service
│   ├── cmd/server/      # Application entry point
│   ├── internal/
│   │   ├── graph/       # GraphQL schema, resolvers, and generated code
│   │   ├── service/     # Business logic layer
│   │   ├── repository/  # Data access layer  
│   │   ├── models/      # Domain models and types
│   │   ├── auth/        # Authentication middleware
│   │   ├── database/    # Database connections and migrations
│   │   ├── cache/       # Redis caching layer
│   │   └── mocks/       # Test mocks
│   ├── Dockerfile       # Backend container configuration
│   └── go.mod          # Go dependencies
├── frontend/            # React TypeScript application
│   ├── src/
│   │   ├── components/  # React components (Login, PetList, Cart, etc.)
│   │   ├── contexts/    # React Context (Auth, Cart state management)
│   │   ├── graphql/     # GraphQL queries and mutations
│   │   ├── config/      # Apollo Client configuration
│   │   └── types/       # TypeScript type definitions
│   ├── Dockerfile       # Frontend container configuration
│   ├── nginx.conf       # Production web server configuration
│   └── package.json     # Node.js dependencies
├── docker-compose.yml   # Multi-service orchestration
└── fullstack-challenge-md.md  # Original requirements
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Git

### Start the Application
```bash
# Clone the repository
git clone https://github.com/fehepe/pet-store
cd pet-store

# Start all services (PostgreSQL, Redis, Backend, Frontend)
docker-compose up -d

# Verify all services are healthy
docker-compose ps
```

### Access the Application
- **🌐 Frontend**: http://localhost:3000
- **🔍 GraphQL Playground**: http://localhost:8080/playground  
- **📡 GraphQL API**: http://localhost:8080/graphql

## 🎯 How to Use

### Customer Experience
1. **Access**: Navigate to http://localhost:3000
2. **Login**: Select a store from the dropdown (no authentication required for demo)
3. **Browse**: View available pets with images and descriptions
4. **Purchase Options**:
   - **Instant Buy**: Click "Buy Now" for immediate purchase
   - **Cart Checkout**: Add multiple pets to cart and checkout together
5. **Validation**: System prevents adding pets already in cart

### API Testing
Use GraphQL Playground at http://localhost:8080/playground:

```graphql
# List all available stores
query {
  listStores {
    id
    name
    createdAt
  }
}

# Browse available pets in a store  
query {
  availablePets(storeID: "123e4567-e89b-12d3-a456-426614174000", pagination: { first: 10 }) {
    edges {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
    }
    totalCount
  }
}
```

## 🏗️ Architecture

### Backend Stack
- **Language**: Go 1.24
- **API**: GraphQL with gqlgen
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Authentication**: Basic HTTP Auth with role-based access
- **Deployment**: Docker containers with multi-stage builds

### Frontend Stack  
- **Framework**: React 18 with TypeScript
- **GraphQL Client**: Apollo Client with caching
- **UI Library**: Material-UI (MUI) v5
- **State Management**: React Context API
- **Form Handling**: React Hook Form
- **Build Tool**: Create React App
- **Web Server**: Nginx (production)

### Data Model
```
Store (1) ←→ (N) Pet
Pet (N) ←→ (N) Order (via OrderPets junction)
```

## 🧪 Testing

### Run Backend Tests
```bash
# Run all tests
cd backend && go test ./...

# Run specific test suites
go test ./internal/service -v
go test ./internal/repository -v  
go test ./internal/validation -v
```

### Test Coverage
- ✅ Store service tests (Create, GetByOwnerID, ListAllStores)
- ✅ Pet service tests with store relationships
- ✅ Order service tests with validation
- ✅ Validation tests for all input types
- ✅ Error handling tests

## 🚀 Deployment

### Production Build
```bash
# Build optimized containers
docker-compose build

# Production deployment
docker-compose -f docker-compose.yml up -d
```

### Environment Configuration
The system uses environment variables configured in `docker-compose.yml`:
- Database connection settings
- Redis configuration  
- Authentication settings
- CORS policies

## 📊 System Status

| Service | Status | Port | Health Check |
|---------|--------|------|--------------|
| Frontend | ✅ Running | 3000 | http://localhost:3000 |
| Backend | ✅ Running | 8080 | GraphQL Playground |
| PostgreSQL | ✅ Running | 5432 | Container health check |
| Redis | ✅ Running | 6379 | Container health check |

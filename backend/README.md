# 🐾 Pet Store GraphQL API

A complete GraphQL API for managing a pet store with merchant and customer functionality. Built with Go, PostgreSQL, Redis, and GraphQL.

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23+ (for local development)

### 1. Start the Application
```bash
# Clone the repository
git clone https://github.com/fehepe/pet-store
cd pet-store/backend

# Start all services (PostgreSQL, Redis, Backend)
docker-compose up -d

# Check status
docker-compose ps
```

### 2. Access the API
- **GraphQL Playground**: [http://localhost:8080/playground](http://localhost:8080/playground)
- **GraphQL Endpoint**: `http://localhost:8080/graphql`
- **Health Check**: `http://localhost:8080/health`

## 🔐 Authentication

The API uses **Basic Authentication** with predefined credentials:

| Role | Username | Password |
|------|----------|----------|
| **Merchant** | `merchant1` | `merchant123` |
| **Customer** | `customer1` | `customer123` |

## 📋 Available Operations

### 🏪 Merchant Operations

#### 1. Create Store
```graphql
mutation CreateStore {
  createStore(input: { name: "My Pet Paradise" }) {
    id
    name
    createdAt
  }
}
```

#### 2. Get My Store
```graphql
query MyStore {
  myStore {
    id
    name
    createdAt
  }
}
```

#### 3. Add Pet
```graphql
mutation CreatePet {
  createPet(input: {
    name: "Fluffy"
    species: Cat
    age: 2
    pictureUrl: "https://example.com/cat.jpg"
    description: "Adorable orange tabby cat"
    breederName: "John Doe"
    breederEmail: "john@example.com"
  }) {
    id
    name
    species
    age
    status
    breederEmail
    createdAt
  }
}
```

#### 4. List All Pets (with filters)
```graphql
query ListPets {
  listPets(
    filter: { status: available }
    pagination: { first: 10, after: "0" }
  ) {
    edges {
      id
      name
      species
      age
      status
      breederEmail
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      endCursor
    }
    totalCount
  }
}
```

#### 5. Get Specific Pet
```graphql
query GetPet {
  getPet(id: "pet-uuid-here") {
    id
    name
    species
    age
    description
    breederName
    breederEmail
    status
    createdAt
  }
}
```

#### 6. View Sales Report
```graphql
query SoldPets {
  soldPets(
    startDate: "2024-01-01T00:00:00Z"
    endDate: "2024-12-31T23:59:59Z"
    pagination: { first: 20 }
  ) {
    edges {
      id
      name
      species
      age
      breederName
      createdAt
    }
    totalCount
  }
}
```

#### 7. View Unsold Pets
```graphql
query UnsoldPets {
  unsoldPets(pagination: { first: 10 }) {
    edges {
      id
      name
      species
      age
      status
    }
    totalCount
  }
}
```

#### 8. Delete Pet
```graphql
mutation DeletePet {
  deletePet(id: "pet-uuid-here")
}
```

### 👥 Customer Operations

#### 1. Browse Available Pets
```graphql
query AvailablePets {
  availablePets(
    storeID: "store-uuid-here"
    pagination: { first: 10 }
  ) {
    edges {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
      # breederEmail is hidden for customers
    }
    pageInfo {
      hasNextPage
      endCursor
    }
    totalCount
  }
}
```

#### 2. Get Store Information
```graphql
query GetStore {
  getStoreByID(id: "store-uuid-here") {
    id
    name
    createdAt
  }
}
```

#### 3. Purchase Single Pet
```graphql
mutation PurchasePet {
  purchasePet(petID: "pet-uuid-here") {
    id
    customerID
    pets {
      id
      name
      species
    }
    totalPets
    createdAt
  }
}
```

#### 4. Purchase Multiple Pets
```graphql
mutation PurchasePets {
  purchasePets(petIDs: ["pet-uuid-1", "pet-uuid-2"]) {
    id
    customerID
    pets {
      id
      name
      species
    }
    totalPets
    createdAt
  }
}
```

## 🎯 Complete Workflow Examples

### Merchant Workflow
```graphql
# 1. First, create a store
mutation { createStore(input: { name: "Fluffy Friends Store" }) { id name } }

# 2. Add some pets
mutation { 
  createPet(input: { 
    name: "Whiskers", species: Cat, age: 1, 
    breederName: "Alice", breederEmail: "alice@breeders.com" 
  }) { id name } 
}

# 3. Check your store and pets
query { myStore { id name } }
query { listPets { edges { id name status } totalCount } }

# 4. View sales in date range
query { 
  soldPets(startDate: "2024-01-01T00:00:00Z", endDate: "2024-12-31T23:59:59Z") { 
    edges { id name } 
    totalCount 
  } 
}
```

### Customer Workflow
```graphql
# 1. Find a store (you'll need the store UUID from merchant)
query { getStoreByID(id: "store-uuid") { id name } }

# 2. Browse available pets
query { 
  availablePets(storeID: "store-uuid", pagination: { first: 5 }) { 
    edges { id name species age description } 
  } 
}

# 3. Purchase a pet
mutation { 
  purchasePet(petID: "pet-uuid") { 
    id customerID totalPets 
    pets { id name } 
  } 
}
```

## 🧪 Testing with Postman

Import the included Postman collection: `Pet-Store-API.postman_collection.json`

**Setup Variables:**
- `base_url`: `http://localhost:8080`

**Authentication:**
- Set to Basic Auth with merchant1/merchant123 or customer1/customer123

## 🛠️ Development

### Local Development Setup
```bash
# Install dependencies
make deps

# Install development tools
make install-tools

# Start only database services
docker-compose up -d postgres redis

# Run locally
make build
./bin/server

# Or use development mode
make dev
```

### Available Commands
```bash
make help              # Show all available commands
make build             # Build the application
make test              # Run all tests
make test-coverage     # Run tests with coverage
make fmt               # Format code
make check             # Run all quality checks
make docker-up         # Start with Docker
make docker-down       # Stop services
make clean             # Clean up
```

### Testing
```bash
# Unit tests
make test-unit

# With coverage
make test-coverage

# Integration tests (requires test DB)
make test-integration
```

## 📊 API Schema Overview

### Types
- **Pet**: Core entity with species (Cat/Dog/Frog), status (available/sold)
- **Store**: Merchant's store with pets
- **Order**: Purchase record with customer and pets
- **PetConnection**: Paginated pet results

### Key Features
- **Pagination**: Cursor-based pagination for all list operations
- **Filtering**: Filter pets by status, date range
- **Authentication**: Role-based access (merchant vs customer)
- **Data Privacy**: Breeder emails hidden from customers
- **Validation**: Input validation for all operations
- **Caching**: Redis caching for performance
- **Encryption**: Breeder emails encrypted at rest

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   GraphQL API   │────│   Service Layer │────│ Repository Layer│
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
    ┌──────────┐         ┌──────────────┐         ┌─────────────┐
    │ Auth     │         │ Redis Cache  │         │ PostgreSQL  │
    │ Middleware│         │             │         │ Database    │
    └──────────┘         └──────────────┘         └─────────────┘
```

## 🔧 Configuration

### Environment Variables
```bash
# Database
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=petstore

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Security
ENCRYPTION_KEY=your-32-character-key-here

# Server
PORT=8080
ENV=development
```

## 🐛 Troubleshooting

### Common Issues

1. **"Connection refused" errors**
   ```bash
   # Ensure services are running
   docker-compose ps
   docker-compose logs backend
   ```

2. **Database connection issues**
   ```bash
   # Reset database
   make db-reset
   ```

3. **Port already in use**
   ```bash
   # Stop existing services
   docker-compose down
   # Or change ports in docker-compose.yml
   ```

### Health Checks
```bash
# Check service health
curl http://localhost:8080/health

# Check database
docker-compose exec postgres pg_isready

# Check Redis
docker-compose exec redis redis-cli ping
```

## 📝 Data Models

### Pet Species
- `Cat`
- `Dog` 
- `Frog`

### Pet Status
- `available` - Can be purchased
- `sold` - Already purchased

### Validation Rules
- Pet names: 1-100 characters
- Pet age: 0-50 years
- Store names: 1-100 characters
- Email format validation
- Maximum 10 pets per order

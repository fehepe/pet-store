# 📡 Pet Store GraphQL API Backend

Complete GraphQL API backend for the Pet Store application with comprehensive merchant and customer functionality.

## 🌐 Endpoint Information

- **GraphQL Endpoint**: `http://localhost:8080/graphql`
- **GraphQL Playground**: `http://localhost:8080/playground`
- **Schema**: Available via introspection
- **Authentication**: Basic HTTP Auth (for non-public endpoints)

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.24+ (for local development)

### Start the Backend
```bash
# Clone the repository
git clone https://github.com/fehepe/pet-store
cd pet-store

# Start all services (PostgreSQL, Redis, Backend)
docker-compose up -d

# Check status
docker-compose ps
```

### Access Points
- **🔍 GraphQL Playground**: http://localhost:8080/playground  
- **📡 GraphQL API**: http://localhost:8080/graphql
- **🏥 Health Check**: http://localhost:8080/health

## 🔓 Public Endpoints

These endpoints are accessible without authentication and are used by the frontend application.

### 📋 List Stores

**Query**: `listStores`
**Description**: Get all available stores for customer selection
**Authentication**: None required

```graphql
query ListStores {
  listStores {
    id
    name
    createdAt
  }
}
```

**Response**:
```json
{
  "data": {
    "listStores": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "Pet Paradise Store",
        "createdAt": "2025-06-26T19:52:03.579789Z"
      }
    ]
  }
}
```

### 🐾 Browse Available Pets

**Query**: `availablePets`
**Description**: Browse available pets in a specific store with pagination
**Authentication**: None required (public browsing)

```graphql
query GetAvailablePets($storeID: UUID!, $pagination: PaginationInput) {
  availablePets(storeID: $storeID, pagination: $pagination) {
    edges {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
      status
      createdAt
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
```

**Variables**:
```json
{
  "storeID": "123e4567-e89b-12d3-a456-426614174000",
  "pagination": {
    "first": 10
  }
}
```

## 🔐 Customer Endpoints

These endpoints require customer authentication using Basic HTTP Auth.

### 🛒 Purchase Single Pet

**Mutation**: `purchasePet`
**Description**: Purchase a single pet immediately
**Authentication**: Customer credentials required

```graphql
mutation PurchasePet($petID: UUID!) {
  purchasePet(petID: $petID) {
    id
    customerID
    pets {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
    }
    totalPets
    createdAt
  }
}
```

### 🛍️ Purchase Multiple Pets

**Mutation**: `purchasePets`
**Description**: Purchase multiple pets in a single transaction
**Authentication**: Customer credentials required

```graphql
mutation PurchasePets($petIDs: [UUID!]!) {
  purchasePets(petIDs: $petIDs) {
    id
    customerID
    pets {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
    }
    totalPets
    createdAt
  }
}
```

## 🏪 Merchant Endpoints

These endpoints require merchant authentication and are used for store management.

### 🏗️ Create Store

**Mutation**: `createStore`
**Description**: Create a new store (one per merchant)
**Authentication**: Merchant credentials required

```graphql
mutation CreateStore($input: CreateStoreInput!) {
  createStore(input: $input) {
    id
    name
    createdAt
  }
}
```

**Variables**:
```json
{
  "input": {
    "name": "My Pet Store"
  }
}
```

### 📋 List My Pets

**Query**: `listPets`
**Description**: List all pets in the merchant's store with filtering
**Authentication**: Merchant credentials required

```graphql
query ListPets($filter: PetFilterInput, $pagination: PaginationInput) {
  listPets(filter: $filter, pagination: $pagination) {
    edges {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
      breederEmail
      status
      createdAt
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
```

### ➕ Add New Pet

**Mutation**: `createPet`
**Description**: Add a new pet to the store inventory
**Authentication**: Merchant credentials required

```graphql
mutation CreatePet($input: CreatePetInput!) {
  createPet(input: $input) {
    id
    name
    species
    age
    pictureUrl
    description
    breederName
    breederEmail
    status
    createdAt
  }
}
```

**Variables**:
```json
{
  "input": {
    "name": "Fluffy",
    "species": "Cat",
    "age": 3,
    "pictureUrl": "https://example.com/fluffy.jpg",
    "description": "A very fluffy and friendly cat",
    "breederName": "Best Cat Breeders",
    "breederEmail": "contact@bestcatbreeders.com"
  }
}
```

### 🔍 Get Pet Details

**Query**: `getPet`
**Description**: Get detailed information about a specific pet
**Authentication**: Merchant credentials required

```graphql
query GetPet($id: UUID!) {
  getPet(id: $id) {
    id
    name
    species
    age
    pictureUrl
    description
    breederName
    breederEmail
    status
    createdAt
  }
}
```

### 🗑️ Delete Pet

**Mutation**: `deletePet`
**Description**: Remove a pet from the store inventory
**Authentication**: Merchant credentials required

```graphql
mutation DeletePet($id: UUID!) {
  deletePet(id: $id)
}
```

### 📊 View Sold Pets

**Query**: `soldPets`
**Description**: View pets sold within a specific date range
**Authentication**: Merchant credentials required

```graphql
query SoldPets($startDate: Time!, $endDate: Time!, $pagination: PaginationInput) {
  soldPets(startDate: $startDate, endDate: $endDate, pagination: $pagination) {
    edges {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
      breederEmail
      status
      createdAt
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
```

### 📦 View Unsold Pets

**Query**: `unsoldPets`
**Description**: View all currently available pets in the store
**Authentication**: Merchant credentials required

```graphql
query UnsoldPets($pagination: PaginationInput) {
  unsoldPets(pagination: $pagination) {
    edges {
      id
      name
      species
      age
      pictureUrl
      description
      breederName
      breederEmail
      status
      createdAt
    }
    pageInfo {
      hasNextPage
      hasPreviousPage
      startCursor
      endCursor
    }
    totalCount
  }
}
```

## 🔧 Types & Schemas

### Pet Species Enum
```graphql
enum PetSpecies {
  Cat
  Dog
  Frog
}
```

### Pet Status Enum
```graphql
enum PetStatus {
  available
  sold
}
```

### Pagination Input
```graphql
input PaginationInput {
  first: Int
  after: String
  last: Int
  before: String
}
```

### Pet Filter Input
```graphql
input PetFilterInput {
  status: PetStatus
  startDate: Time
  endDate: Time
}
```

### Create Pet Input
```graphql
input CreatePetInput {
  name: String!
  species: PetSpecies!
  age: Int!
  pictureUrl: String
  description: String
  breederName: String!
  breederEmail: String!
}
```

### Create Store Input
```graphql
input CreateStoreInput {
  name: String!
}
```

## 🔑 Authentication

The API uses Basic HTTP Authentication with role-based access:

### Merchant Authentication
```bash
Authorization: Basic <base64(username:password)>
```

**Example**:
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(echo -n 'merchant1:merchant123' | base64)" \
  -d '{"query": "query { myStore { id name } }"}'
```

### Customer Authentication
```bash
Authorization: Basic <base64(username:password)>
```

**Example**:
```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Basic $(echo -n 'customer1:customer123' | base64)" \
  -d '{"query": "mutation { purchasePet(petID: \"pet-id\") { id } }"}'
```

## 📝 Error Handling

### GraphQL Errors
Errors are returned in the standard GraphQL format:

```json
{
  "errors": [
    {
      "message": "Pet not found",
      "path": ["getPet"],
      "extensions": {
        "code": "NOT_FOUND"
      }
    }
  ],
  "data": null
}
```

### Common Error Types
- **Authentication Required**: `401 Unauthorized`
- **Forbidden**: `403 Forbidden` (wrong role)
- **Validation Error**: GraphQL validation errors
- **Not Found**: Resource not found errors
- **Conflict**: Resource already exists

## 🧪 Testing with Postman

1. **Import Collection**: Import `Pet-Store-API.postman_collection.json`
2. **Import Environment**: Import `Pet-Store-Environment.postman_environment.json`
3. **Set Variables**: Update environment variables as needed
4. **Run Requests**: Execute requests in logical order

### Test Data
The system includes seed data with:
- 1 store: "Pet Paradise Store"
- 11 pets with various species (cats, dogs, frogs)
- Mix of pets with and without pictures

## 🔄 Recent API Changes

### Removed Endpoints
- ❌ `getStoreByID` - Replaced with `listStores` for better efficiency

### Enhanced Endpoints
- ✅ `listStores` - Now public endpoint for store selection
- ✅ `availablePets` - Public browsing without authentication

### Added Features
- ✅ Public access to store listings and pet browsing
- ✅ Enhanced error messages and validation
- ✅ Better pagination support

## 🛠️ Development

### Local Development Setup
```bash
# Install dependencies
go mod download

# Start only database services
docker-compose up -d postgres redis

# Set environment variables
export DB_HOST=localhost
export REDIS_HOST=localhost

# Run the server
go run cmd/server/main.go
```

### Available Commands
```bash
# Run tests
go test ./...

# Run specific test suites
go test ./internal/service -v
go test ./internal/repository -v

# Generate GraphQL code
go generate ./internal/graph

# Build the application
go build -o bin/server cmd/server/main.go

# Force re-run migrations (if needed)
docker-compose exec backend /app/main migrate
```

### Testing
```bash
# Unit tests
go test ./internal/service/...

# Integration tests
go test ./internal/repository/...

# With coverage
go test -cover ./...
```

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

### Backend Stack
- **Language**: Go 1.24
- **API**: GraphQL with gqlgen
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Authentication**: Basic HTTP Auth with role-based access
- **Deployment**: Docker containers with multi-stage builds

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

## 🚀 Development Tools

### GraphQL Playground
Access the interactive GraphQL Playground at:
`http://localhost:8080/playground`

Features:
- Schema documentation
- Query/mutation autocompletion
- Real-time query execution
- Schema introspection

### Schema Introspection
Get the complete schema programmatically:

```graphql
query IntrospectionQuery {
  __schema {
    types {
      name
      kind
      description
      fields {
        name
        type {
          name
          kind
        }
      }
    }
  }
}
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
   # Check database status
   docker-compose exec postgres pg_isready
   ```

3. **Port already in use**
   ```bash
   # Stop existing services
   docker-compose down
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

## 📊 System Status

| Service | Status | Port | Health Check |
|---------|--------|------|--------------|
| Backend API | ✅ Running | 8080 | http://localhost:8080/health |
| PostgreSQL | ✅ Running | 5432 | Container health check |
| Redis | ✅ Running | 6379 | Container health check |
| GraphQL Playground | ✅ Running | 8080 | http://localhost:8080/playground |

## 📞 Support

For API issues or questions:
1. Check the GraphQL Playground for schema documentation
2. Review error messages in the response
3. Check authentication credentials and permissions
4. Ensure all required variables are provided
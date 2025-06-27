# Pet Store API Backend

GraphQL API backend for the Pet Store application.

## Quick Start

```bash
docker-compose up -d
```

**Access:**
- GraphQL Playground: http://localhost:8080/playground
- GraphQL API: http://localhost:8080/graphql

## API Endpoints

### Public (No Auth Required)

**List Stores**
```graphql
{ listStores { id name createdAt } }
```

**Browse Pets**
```graphql
{ 
  availablePets(storeID: "123e4567-e89b-12d3-a456-426614174000") {
    edges { id name species age pictureUrl description }
    totalCount
  }
}
```

### Customer (Auth Required)

**Purchase Pet**
```graphql
mutation { 
  purchasePet(petID: "pet-id") { 
    id customerID totalPets 
  } 
}
```

**Purchase Multiple Pets**
```graphql
mutation { 
  purchasePets(petIDs: ["pet-id-1", "pet-id-2"]) { 
    id customerID totalPets 
  } 
}
```

### Merchant (Auth Required)

**Create Store**
```graphql
mutation { 
  createStore(input: {name: "My Pet Store"}) { 
    id name 
  } 
}
```

**Add Pet**
```graphql
mutation { 
  createPet(input: {
    name: "Fluffy"
    species: Cat
    age: 3
    breederName: "Best Breeders"
    breederEmail: "contact@breeders.com"
  }) { 
    id name species 
  } 
}
```

**List My Pets**
```graphql
{ 
  listPets { 
    edges { id name species status } 
  } 
}
```

## Authentication

Use Basic HTTP Auth:
- **Customer**: `customer1:customer123`
- **Merchant**: `merchant1:merchant123`

```bash
curl -H "Authorization: Basic $(echo -n 'customer1:customer123' | base64)" \
     -H "Content-Type: application/json" \
     -d '{"query": "mutation { purchasePet(petID: \"pet-id\") { id } }"}' \
     http://localhost:8080/graphql
```

## Development

```bash
# Local development
go mod download
docker-compose up -d postgres redis
export DB_HOST=localhost REDIS_HOST=localhost
go run cmd/server/main.go

# Tests
go test ./...

# Generate GraphQL code
go generate ./internal/graph
```

## Tech Stack

- Go 1.24
- GraphQL (gqlgen)
- PostgreSQL 15
- Redis 7
- Docker
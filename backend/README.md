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

## Postman Screenshots
![Screenshot 2025-06-27 at 1 15 03 AM](https://github.com/user-attachments/assets/2c2e0d84-5f0c-4414-8d08-b7c2df215fad)
![Screenshot 2025-06-27 at 1 15 15 AM](https://github.com/user-attachments/assets/373e0c9c-7911-454d-91d0-02a6bff7b6ca)
![Screenshot 2025-06-27 at 1 15 35 AM](https://github.com/user-attachments/assets/7e8a1afa-cbd1-4281-8ec9-b69ad6bf4c2c)
![Screenshot 2025-06-27 at 1 16 45 AM](https://github.com/user-attachments/assets/5c69c256-388a-452d-9c47-4122c1773c77)
![Screenshot 2025-06-27 at 1 16 58 AM](https://github.com/user-attachments/assets/381a25aa-7025-4015-975e-980480eadcbf)
![Screenshot 2025-06-27 at 1 17 12 AM](https://github.com/user-attachments/assets/db64341b-f1a8-4ae4-9d74-b7afc9e41f7a)
![Screenshot 2025-06-27 at 1 17 33 AM](https://github.com/user-attachments/assets/fb5cb778-8042-4b2b-965e-ade23264b448)
![Screenshot 2025-06-27 at 1 17 53 AM](https://github.com/user-attachments/assets/518cd5ac-0d0c-4fd7-b178-76b955ee592f)
![Screenshot 2025-06-27 at 1 18 05 AM](https://github.com/user-attachments/assets/eb47c00d-210a-407d-a4f5-f67ecdec82d7)
![Screenshot 2025-06-27 at 1 18 16 AM](https://github.com/user-attachments/assets/a35e120e-5598-4504-9091-8ed43841d15a)
![Screenshot 2025-06-27 at 1 18 26 AM](https://github.com/user-attachments/assets/4fb1c610-a43c-48e1-8ef6-d56769585634)

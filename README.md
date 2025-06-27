# 🐾 Pet Store

Full-stack pet store application with GraphQL API backend (Go) and React TypeScript frontend.

## Features

- Browse pets by store
- Shopping cart with instant purchase
- Real-time inventory updates
- Responsive Material-UI design
- GraphQL API with authentication

## Quick Start

```bash
git clone https://github.com/fehepe/pet-store
cd pet-store
docker-compose up -d
```

**Access:**
- Frontend: http://localhost:3000
- GraphQL Playground: http://localhost:8080/playground

## Usage

1. Open http://localhost:3000
2. Select a store from dropdown
3. Browse pets and add to cart
4. Purchase instantly or checkout multiple pets

## Tech Stack

**Backend:** Go, GraphQL, PostgreSQL, Redis  
**Frontend:** React, TypeScript, Material-UI, Apollo Client  
**Deployment:** Docker Compose

## API Testing

Test GraphQL queries at http://localhost:8080/playground:

```graphql
# List stores
{ listStores { id name } }

# Browse pets
{ availablePets(storeID: "123e4567-e89b-12d3-a456-426614174000") 
  { edges { id name species age } } }
```

## Development

```bash
# Backend tests
cd backend && go test ./...

# Frontend development
cd frontend && npm start
```

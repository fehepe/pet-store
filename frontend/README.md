# Pet Store Frontend

React TypeScript frontend for the Pet Store application.

## Features

- Browse pets by store with responsive grid layout
- Shopping cart with instant purchase options
- Material-UI design system
- Real-time inventory updates
- Default images for pets without pictures

## Quick Start

```bash
# With Docker (recommended)
docker-compose up -d

# Local development
npm install
npm start
```

**Access:** http://localhost:3000

## Usage

1. Select a store from the dropdown
2. Browse available pets
3. Either:
   - Click "Buy Now" for instant purchase
   - Add pets to cart and checkout multiple pets

## Tech Stack

- React 18 with TypeScript
- Material-UI (MUI) v5
- Apollo Client (GraphQL)
- React Context API (state management)
- Docker & Nginx (deployment)

## Development

```bash
# Install dependencies
npm install

# Start development server
npm start

# Run tests
npm test

# Build for production
npm run build
```

## Environment Variables

Create `.env` file:
```bash
REACT_APP_GRAPHQL_ENDPOINT=http://localhost:8080/graphql
```

## Project Structure

```
src/
├── components/     # React components
├── contexts/       # Authentication and cart state
├── graphql/        # GraphQL queries and mutations
├── config/         # Apollo Client setup
└── types/          # TypeScript definitions
```
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

## Sc![Screenshot 2025-06-26 at 11 23 33 AM](https://github.com/user-attachments/assets/475a4f98-ac5d-4d41-9c78-364e464ca64c)
![Screenshot 2025-06-26 at 10 46 10 PM](https://github.com/user-attachments/assets/2f0a93dc-f954-4c33-b53e-58d4137972d5)
![Screenshot 2025-06-26 at 10 46 23 PM](https://github.com/user-attachments/assets/e9d9725d-c647-462d-a8c0-4de930a1b20c)
![Screenshot 2025-06-26 at 10 46 41 PM](https://github.com/user-attachments/assets/6718b48e-2f8d-4805-a3e4-42291bfcba31)
![Screenshot 2025-06-27 at 12 42 47 AM](https://github.com/user-attachments/assets/bd542816-8bf0-4bfd-8c7e-ca147533d600)
reenshots


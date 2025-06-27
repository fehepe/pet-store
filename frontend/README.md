# 🌐 Pet Store Frontend

A modern React TypeScript application providing a customer-facing interface for the Pet Store system.

## 🎯 Overview

This frontend application provides an intuitive interface for customers to browse and purchase pets from different stores. Built with React 18, TypeScript, and Material-UI, it offers a responsive and professional user experience.

## ✨ Features

### 🏪 Store Management
- **Store Selection**: User-friendly dropdown populated from backend API
- **Dynamic Loading**: Stores fetched in real-time from `listStores` endpoint
- **No Manual Entry**: Eliminated error-prone UUID input

### 🐾 Pet Browsing  
- **Grid Display**: Responsive grid layout showing available pets
- **Rich Information**: Name, species, age, description, and breeder details
- **Image Handling**: 
  - Custom pet images when available
  - Species-specific default images (cats, dogs, frogs)
  - Fallback images from Unsplash

### 🛒 Shopping Experience
- **Instant Purchase**: One-click "Buy Now" functionality
- **Shopping Cart**: Add multiple pets for bulk checkout
- **Smart Validation**: 
  - Prevents adding pets already in cart
  - Visual feedback with disabled buttons and tooltips
  - "In Cart" status indication

### 🎨 User Interface
- **Material-UI Components**: Professional, accessible design system
- **Responsive Design**: Works on desktop, tablet, and mobile
- **Dark/Light Theme**: Consistent theming across components
- **Loading States**: Proper loading indicators and error handling
- **Toast Notifications**: Success and error feedback

## 🏗️ Architecture

### Technology Stack
```
React 18.2.0          # UI framework
TypeScript 4.9.5      # Type safety
Material-UI 5.14.15   # Component library
Apollo Client 3.8.6   # GraphQL client
React Hook Form 7.47.0 # Form management
React Router 6.17.0   # Navigation (if routing added)
```

### Project Structure
```
src/
├── components/          # React components
│   ├── Layout.tsx      # Main layout wrapper
│   ├── Login.tsx       # Store selection and authentication  
│   ├── PetList.tsx     # Pet grid display with pagination
│   ├── PetCard.tsx     # Individual pet card component
│   ├── Cart.tsx        # Shopping cart sidebar
│   └── ProtectedRoute.tsx # Authentication guard
├── contexts/           # React Context providers
│   ├── AuthContext.tsx # Authentication state management
│   └── CartContext.tsx # Shopping cart state management
├── graphql/           # GraphQL queries and mutations
│   └── queries.ts     # All GraphQL operations
├── config/            # Configuration
│   └── apollo.ts      # Apollo Client setup
├── types/             # TypeScript definitions
│   └── index.ts       # Shared type definitions
├── theme.ts           # Material-UI theme configuration
├── App.tsx            # Main application component
└── index.tsx          # Application entry point
```

### State Management
The application uses React Context API for state management:

**AuthContext** - Manages authentication state:
```typescript
interface AuthContextType {
  user: string | null;
  storeId: string | null;
  login: (username: string, storeId: string) => void;
  logout: () => void;
  loading: boolean;
}
```

**CartContext** - Manages shopping cart:
```typescript
interface CartContextType {
  cartItems: Pet[];
  addToCart: (pet: Pet) => void;
  removeFromCart: (petId: string) => void;
  clearCart: () => void;
  isInCart: (petId: string) => boolean;
}
```

## 🔌 API Integration

### GraphQL Operations

**Store Selection:**
```typescript
const LIST_STORES = gql`
  query ListStores {
    listStores {
      id
      name
      createdAt
    }
  }
`;
```

**Pet Browsing:**
```typescript
const GET_AVAILABLE_PETS = gql`
  query GetAvailablePets($storeID: UUID!, $pagination: PaginationInput) {
    availablePets(storeID: $storeID, pagination: $pagination) {
      edges { ...PetFields }
      pageInfo { hasNextPage endCursor }
      totalCount
    }
  }
`;
```

**Pet Purchase:**
```typescript
const PURCHASE_PET = gql`
  mutation PurchasePet($petID: UUID!) {
    purchasePet(petID: $petID) {
      id
      pets { ...PetFields }
      totalPets
    }
  }
`;
```

### Apollo Client Configuration
```typescript
const client = new ApolloClient({
  uri: process.env.REACT_APP_GRAPHQL_ENDPOINT || 'http://localhost:8080/graphql',
  cache: new InMemoryCache({
    typePolicies: {
      PetConnection: {
        fields: {
          edges: {
            merge: (existing = [], incoming = []) => [...existing, ...incoming]
          }
        }
      }
    }
  }),
  defaultOptions: {
    watchQuery: { errorPolicy: 'all' },
    query: { errorPolicy: 'all' }
  }
});
```

## 🚀 Development

### Prerequisites
- Node.js 18+
- npm or yarn

### Local Development
```bash
# Install dependencies
npm install

# Start development server
npm start

# Open browser to http://localhost:3000
```

### Available Scripts

| Command | Description |
|---------|-------------|
| `npm start` | Start development server with hot reload |
| `npm build` | Create production build |
| `npm test` | Run test suite |
| `npm eject` | Eject from Create React App (one-way) |

### Environment Variables
Create `.env` file in frontend root:
```bash
REACT_APP_GRAPHQL_ENDPOINT=http://localhost:8080/graphql
REACT_APP_API_URL=http://localhost:8080
```

## 🐳 Docker Deployment

### Development
```bash
# Run with backend services
docker-compose up -d
```

### Production Build
```dockerfile
FROM node:18-alpine as build
WORKDIR /app
COPY package*.json ./
RUN npm ci --silent
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 3000
CMD ["nginx", "-g", "daemon off;"]
```

## 🎨 Styling & Theming

### Material-UI Theme
```typescript
const theme = createTheme({
  palette: {
    primary: { main: '#1976d2' },
    secondary: { main: '#dc004e' }
  },
  typography: {
    fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif'
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: { textTransform: 'none' }
      }
    }
  }
});
```

## 📱 User Experience

### Navigation Flow
1. **Landing** → Login page with store selection
2. **Store Selection** → Dropdown populated from API  
3. **Pet Browsing** → Grid of available pets with pagination
4. **Purchase Options** → Instant buy or add to cart
5. **Checkout** → Review and confirm purchase

## 🧪 Testing

### Test Structure
```bash
# Run all tests
npm test

# Run with coverage
npm test -- --coverage

# Run specific tests  
npm test -- --testNamePattern="PetCard"
```

### Testing Libraries
- **Jest**: Test runner and assertions
- **React Testing Library**: Component testing
- **Apollo MockedProvider**: GraphQL mocking
- **MSW**: API mocking for integration tests

## 🚀 Deployment

### Nginx Configuration
```nginx
server {
    listen 3000;
    location / {
        root /usr/share/nginx/html;
        index index.html index.htm;
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass http://backend:8080;
    }
}
```

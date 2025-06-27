import { gql } from '@apollo/client';

export const PET_FRAGMENT = gql`
  fragment PetFields on Pet {
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
`;

export const GET_AVAILABLE_PETS = gql`
  ${PET_FRAGMENT}
  query GetAvailablePets($storeID: UUID!, $pagination: PaginationInput) {
    availablePets(storeID: $storeID, pagination: $pagination) {
      edges {
        ...PetFields
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
`;

export const LIST_STORES = gql`
  query ListStores {
    listStores {
      id
      name
      createdAt
    }
  }
`;

export const PURCHASE_PET = gql`
  ${PET_FRAGMENT}
  mutation PurchasePet($petID: UUID!) {
    purchasePet(petID: $petID) {
      id
      customerID
      pets {
        ...PetFields
      }
      totalPets
      createdAt
    }
  }
`;

export const PURCHASE_PETS = gql`
  ${PET_FRAGMENT}
  mutation PurchasePets($petIDs: [UUID!]!) {
    purchasePets(petIDs: $petIDs) {
      id
      customerID
      pets {
        ...PetFields
      }
      totalPets
      createdAt
    }
  }
`;
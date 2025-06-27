export enum PetSpecies {
  Cat = 'Cat',
  Dog = 'Dog',
  Frog = 'Frog',
}

export enum PetStatus {
  available = 'available',
  sold = 'sold',
}

export interface Pet {
  id: string;
  name: string;
  species: PetSpecies;
  age: number;
  pictureUrl?: string;
  description?: string;
  breederName: string;
  breederEmail: string;
  status: PetStatus;
  createdAt: string;
}

export interface Store {
  id: string;
  name: string;
  createdAt: string;
}

export interface Order {
  id: string;
  customerID: string;
  pets: Pet[];
  totalPets: number;
  createdAt: string;
}

export interface PageInfo {
  hasNextPage: boolean;
  hasPreviousPage: boolean;
  startCursor?: string;
  endCursor?: string;
}

export interface PetConnection {
  edges: Pet[];
  pageInfo: PageInfo;
  totalCount: number;
}

export interface CartItem {
  pet: Pet;
  addedAt: number;
}
export enum PetSpecies {
  Cat = 'Cat',
  Dog = 'Dog',
  Frog = 'Frog',
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
  status: string;
  createdAt: string;
}

export interface CartItem {
  pet: Pet;
  addedAt: number;
}
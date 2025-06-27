import React from 'react';
import PetsIcon from '@mui/icons-material/Pets';

/**
 * Returns the appropriate icon for a pet species
 */
export const getSpeciesIcon = (species: string): React.ReactElement => {
  switch (species.toLowerCase()) {
    case 'cat':
      return React.createElement(PetsIcon);
    case 'dog':
      return React.createElement(PetsIcon);
    case 'frog':
      return React.createElement(PetsIcon);
    default:
      return React.createElement(PetsIcon);
  }
};

/**
 * Returns a default image URL based on pet species
 */
export const getDefaultImage = (species: string): string => {
  const speciesImages = {
    cat: 'https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=300&h=200&fit=crop&crop=center',
    dog: 'https://images.unsplash.com/photo-1552053831-71594a27632d?w=300&h=200&fit=crop&crop=center',
    frog: 'https://images.ctfassets.net/cnu0m8re1exe/4txgybYHtUH0z6Dy9IIFGH/e9f4612ef512ae7ad8a580a39557e4d2/Glass_Frog.jpg?fm=jpg&fl=progressive&w=660&h=433&fit=fill',
  };
  
  return speciesImages[species.toLowerCase() as keyof typeof speciesImages] || speciesImages.cat;
};
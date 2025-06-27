import React from 'react';
import {
  Card,
  CardMedia,
  CardContent,
  CardActions,
  Typography,
  Button,
  Chip,
  Box,
  Tooltip,
} from '@mui/material';
import {
  ShoppingCart,
  Pets,
  AttachMoney,
} from '@mui/icons-material';
import { Pet, PetSpecies } from '../types';
import { useCart } from '../contexts/CartContext';

interface PetCardProps {
  pet: Pet;
  onPurchase: (pet: Pet) => void;
}

const getSpeciesIcon = (species: PetSpecies) => {
  switch (species) {
    case PetSpecies.Cat:
      return '🐱';
    case PetSpecies.Dog:
      return '🐕';
    case PetSpecies.Frog:
      return '🐸';
    default:
      return '🐾';
  }
};

const getSpeciesColor = (species: PetSpecies) => {
  switch (species) {
    case PetSpecies.Cat:
      return 'primary';
    case PetSpecies.Dog:
      return 'secondary';
    case PetSpecies.Frog:
      return 'success';
    default:
      return 'default';
  }
};

const getDefaultImage = (species: PetSpecies) => {
  switch (species) {
    case PetSpecies.Cat:
      return 'https://images.unsplash.com/photo-1514888286974-6c03e2ca1dba?w=300&h=200&fit=crop&crop=center';
    case PetSpecies.Dog:
      return 'https://images.unsplash.com/photo-1552053831-71594a27632d?w=300&h=200&fit=crop&crop=center';
    case PetSpecies.Frog:
      return 'https://images.unsplash.com/photo-1551991545-af0bb7fad89b?w=300&h=200&fit=crop&crop=center';
    default:
      return 'https://images.unsplash.com/photo-1601758228041-f3b2795255f1?w=300&h=200&fit=crop&crop=center';
  }
};

export const PetCard: React.FC<PetCardProps> = ({ pet, onPurchase }) => {
  const { addToCart, removeFromCart, isInCart } = useCart();
  const inCart = isInCart(pet.id);

  const handleCartToggle = () => {
    if (inCart) {
      removeFromCart(pet.id);
    } else {
      addToCart(pet);
    }
  };

  return (
    <Card
      sx={{
        height: '100%',
        display: 'flex',
        flexDirection: 'column',
        transition: 'transform 0.2s',
        '&:hover': {
          transform: 'translateY(-4px)',
          boxShadow: 3,
        },
      }}
    >
      <CardMedia
        component="img"
        height="200"
        image={pet.pictureUrl || getDefaultImage(pet.species)}
        alt={pet.name}
        sx={{ objectFit: 'cover' }}
      />
      <CardContent sx={{ flexGrow: 1, pb: 1 }}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
          <Typography variant="h5" component="h2" fontWeight="bold">
            {pet.name}
          </Typography>
          <Chip
            icon={<Pets />}
            label={`${getSpeciesIcon(pet.species)} ${pet.species}`}
            color={getSpeciesColor(pet.species)}
            size="small"
          />
        </Box>

        <Typography variant="body2" color="text.secondary" gutterBottom>
          Age: {pet.age} {pet.age === 1 ? 'year' : 'years'}
        </Typography>

        {pet.description && (
          <Typography
            variant="body2"
            sx={{
              mt: 1,
              mb: 1,
              display: '-webkit-box',
              WebkitLineClamp: 3,
              WebkitBoxOrient: 'vertical',
              overflow: 'hidden',
            }}
          >
            {pet.description}
          </Typography>
        )}

        <Box mt={2}>
          <Typography variant="caption" display="block" color="text.secondary">
            Breeder: {pet.breederName}
          </Typography>
          <Tooltip title={pet.breederEmail}>
            <Typography
              variant="caption"
              display="block"
              color="text.secondary"
              sx={{
                overflow: 'hidden',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
              }}
            >
              Contact: {pet.breederEmail}
            </Typography>
          </Tooltip>
        </Box>
      </CardContent>

      <CardActions sx={{ p: 2, pt: 0 }}>
        <Tooltip 
          title={inCart ? 'Remove from cart to buy individually' : 'Purchase this pet immediately'}
          arrow
        >
          <span>
            <Button
              variant="contained"
              color="primary"
              fullWidth
              startIcon={<AttachMoney />}
              onClick={() => onPurchase(pet)}
              disabled={inCart}
              sx={{ mr: 1 }}
            >
              {inCart ? 'In Cart' : 'Buy Now'}
            </Button>
          </span>
        </Tooltip>
        <Button
          variant={inCart ? 'outlined' : 'contained'}
          color={inCart ? 'error' : 'secondary'}
          fullWidth
          startIcon={<ShoppingCart />}
          onClick={handleCartToggle}
        >
          {inCart ? 'Remove' : 'Add to Cart'}
        </Button>
      </CardActions>
    </Card>
  );
};
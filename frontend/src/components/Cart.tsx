import React, { useState } from 'react';
import { useMutation } from '@apollo/client';
import {
  Drawer,
  Box,
  Typography,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  IconButton,
  Button,
  Divider,
  Alert,
  Chip,
  ListItemSecondaryAction,
} from '@mui/material';
import {
  Close,
  Delete,
  ShoppingCart,
  Pets,
} from '@mui/icons-material';
import { useCart } from '../contexts/CartContext';
import { PURCHASE_PETS, GET_AVAILABLE_PETS } from '../graphql/queries';
import { PetSpecies } from '../types';
import { useAuth } from '../contexts/AuthContext';

interface CartProps {
  open: boolean;
  onClose: () => void;
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

export const Cart: React.FC<CartProps> = ({ open, onClose }) => {
  const { cartItems, removeFromCart, clearCart, getCartPetIds } = useCart();
  const { storeId } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<boolean>(false);

  const [purchasePets, { loading }] = useMutation(PURCHASE_PETS, {
    refetchQueries: [
      {
        query: GET_AVAILABLE_PETS,
        variables: {
          storeID: storeId,
          pagination: { first: 12 },
        },
      },
    ],
    onCompleted: () => {
      setSuccess(true);
      clearCart();
      setTimeout(() => {
        onClose();
        setSuccess(false);
      }, 2000);
    },
    onError: (error) => {
      if (error.message.includes('already been sold')) {
        // Parse the error to find which pets are unavailable
        const unavailablePets = cartItems
          .filter((item) => error.message.includes(item.pet.name))
          .map((item) => item.pet.name);
        
        if (unavailablePets.length > 0) {
          setError(
            `The following pets are no longer available: ${unavailablePets.join(', ')}. Please remove them from your cart.`
          );
        } else {
          setError('Some pets in your cart are no longer available.');
        }
      } else {
        setError(error.message);
      }
    },
  });

  const handleCheckout = async () => {
    setError(null);
    const petIds = getCartPetIds();
    if (petIds.length === 0) return;

    try {
      await purchasePets({
        variables: { petIDs: petIds },
      });
    } catch (err) {
      // Error handled in onError callback
    }
  };

  const total = cartItems.length;

  return (
    <Drawer
      anchor="right"
      open={open}
      onClose={onClose}
      PaperProps={{
        sx: { width: { xs: '100%', sm: 400 } },
      }}
    >
      <Box sx={{ p: 2 }}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h5" display="flex" alignItems="center" gap={1}>
            <ShoppingCart /> Shopping Cart
          </Typography>
          <IconButton onClick={onClose}>
            <Close />
          </IconButton>
        </Box>

        {success && (
          <Alert severity="success" sx={{ mb: 2 }}>
            Purchase successful! Thank you for your order.
          </Alert>
        )}

        {error && (
          <Alert severity="error" sx={{ mb: 2 }} onClose={() => setError(null)}>
            {error}
          </Alert>
        )}

        {cartItems.length === 0 ? (
          <Box textAlign="center" py={4}>
            <Typography variant="body1" color="text.secondary">
              Your cart is empty
            </Typography>
          </Box>
        ) : (
          <>
            <List sx={{ flexGrow: 1, overflow: 'auto' }}>
              {cartItems.map((item) => (
                <ListItem key={item.pet.id} sx={{ px: 0 }}>
                  <ListItemAvatar>
                    <Avatar
                      src={item.pet.pictureUrl || getDefaultImage(item.pet.species)}
                      alt={item.pet.name}
                      sx={{ width: 56, height: 56 }}
                    >
                      <Pets />
                    </Avatar>
                  </ListItemAvatar>
                  <ListItemText
                    primary={
                      <Box display="flex" alignItems="center" gap={1}>
                        <Typography variant="subtitle1">{item.pet.name}</Typography>
                        <Typography variant="caption">
                          {getSpeciesIcon(item.pet.species)}
                        </Typography>
                      </Box>
                    }
                    secondary={
                      <Box>
                        <Typography variant="body2" color="text.secondary">
                          {item.pet.age} {item.pet.age === 1 ? 'year' : 'years'} old
                        </Typography>
                        <Chip
                          label={item.pet.species}
                          size="small"
                          sx={{ mt: 0.5 }}
                        />
                      </Box>
                    }
                  />
                  <ListItemSecondaryAction>
                    <IconButton
                      edge="end"
                      onClick={() => removeFromCart(item.pet.id)}
                      color="error"
                    >
                      <Delete />
                    </IconButton>
                  </ListItemSecondaryAction>
                </ListItem>
              ))}
            </List>

            <Divider sx={{ my: 2 }} />

            <Box>
              <Typography variant="h6" gutterBottom>
                Total: {total} {total === 1 ? 'pet' : 'pets'}
              </Typography>

              <Button
                variant="contained"
                color="primary"
                fullWidth
                size="large"
                onClick={handleCheckout}
                disabled={loading || cartItems.length === 0}
                sx={{ mt: 2 }}
              >
                {loading ? 'Processing...' : 'Checkout'}
              </Button>

              <Button
                variant="outlined"
                fullWidth
                size="small"
                onClick={clearCart}
                sx={{ mt: 1 }}
              >
                Clear Cart
              </Button>
            </Box>
          </>
        )}
      </Box>
    </Drawer>
  );
};
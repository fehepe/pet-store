import React, { useState, useCallback } from 'react';
import { useQuery, useMutation } from '@apollo/client';
import {
  Container,
  Typography,
  Box,
  CircularProgress,
  Alert,
  Button,
  Snackbar,
  Paper,
  Chip,
} from '@mui/material';
import { Refresh, Store } from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import { Pet } from '../types';
import { PetCard } from './PetCard';
import { GET_AVAILABLE_PETS, PURCHASE_PET, LIST_STORES } from '../graphql/queries';

const PETS_PER_PAGE = 12;

export const PetList: React.FC = () => {
  const { storeId } = useAuth();
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const { data: storesData } = useQuery(LIST_STORES);
  const currentStore = storesData?.listStores?.find((store: any) => store.id === storeId);

  const {
    loading,
    error,
    data,
    fetchMore,
    refetch,
  } = useQuery(GET_AVAILABLE_PETS, {
    variables: {
      storeID: storeId,
      pagination: {
        first: PETS_PER_PAGE,
      },
    },
    skip: !storeId,
    notifyOnNetworkStatusChange: true,
  });

  const [purchasePet] = useMutation(PURCHASE_PET, {
    onCompleted: (data) => {
      setSuccessMessage(`Successfully purchased ${data.purchasePet.pets[0].name}!`);
      refetch();
    },
    onError: (error) => {
      if (error.message.includes('already been sold')) {
        setErrorMessage(`This pet is no longer available for purchase.`);
      } else {
        setErrorMessage(error.message);
      }
    },
  });

  const handlePurchase = useCallback(async (pet: Pet) => {
    try {
      await purchasePet({
        variables: { petID: pet.id },
      });
    } catch (err) {
      // Error handled in onError callback
    }
  }, [purchasePet]);

  const handleLoadMore = () => {
    if (data?.availablePets.pageInfo.hasNextPage) {
      fetchMore({
        variables: {
          pagination: {
            first: PETS_PER_PAGE,
            after: data.availablePets.pageInfo.endCursor,
          },
        },
      });
    }
  };

  const handleRefresh = () => {
    refetch();
  };

  if (!storeId) {
    return (
      <Container maxWidth="lg">
        <Alert severity="error">No store selected</Alert>
      </Container>
    );
  }

  if (loading && !data) {
    return (
      <Container maxWidth="lg">
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="60vh">
          <CircularProgress />
        </Box>
      </Container>
    );
  }

  if (error) {
    return (
      <Container maxWidth="lg">
        <Alert severity="error">Error loading pets: {error.message}</Alert>
      </Container>
    );
  }

  const pets = data?.availablePets.edges || [];
  const hasMore = data?.availablePets.pageInfo.hasNextPage || false;

  return (
    <Container maxWidth="lg" sx={{ py: 4 }}>
      <Paper elevation={0} sx={{ p: 3, mb: 4, bgcolor: 'background.default' }}>
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
          <Box>
            <Typography variant="h4" component="h1" gutterBottom>
              Available Pets
            </Typography>
            {currentStore && (
              <Box display="flex" alignItems="center" gap={1}>
                <Store fontSize="small" />
                <Typography variant="h6" color="text.secondary">
                  {currentStore.name}
                </Typography>
                <Chip
                  label={`${data?.availablePets.totalCount || 0} pets available`}
                  color="primary"
                  size="small"
                />
              </Box>
            )}
          </Box>
          <Button
            variant="outlined"
            startIcon={<Refresh />}
            onClick={handleRefresh}
            disabled={loading}
          >
            Refresh
          </Button>
        </Box>
      </Paper>

      {pets.length === 0 ? (
        <Paper sx={{ p: 4, textAlign: 'center' }}>
          <Typography variant="h6" color="text.secondary">
            No pets available at this store right now.
          </Typography>
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            Please check back later or refresh the page.
          </Typography>
        </Paper>
      ) : (
        <>
          <Box
            display="grid"
            gridTemplateColumns={{
              xs: '1fr',
              sm: 'repeat(2, 1fr)',
              md: 'repeat(3, 1fr)',
              lg: 'repeat(4, 1fr)',
            }}
            gap={3}
          >
            {pets.map((pet: Pet) => (
              <PetCard key={pet.id} pet={pet} onPurchase={handlePurchase} />
            ))}
          </Box>

          {hasMore && (
            <Box display="flex" justifyContent="center" mt={4}>
              <Button
                variant="contained"
                onClick={handleLoadMore}
                disabled={loading}
              >
                {loading ? <CircularProgress size={24} /> : 'Load More Pets'}
              </Button>
            </Box>
          )}
        </>
      )}

      <Snackbar
        open={!!successMessage}
        autoHideDuration={4000}
        onClose={() => setSuccessMessage(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert severity="success" onClose={() => setSuccessMessage(null)}>
          {successMessage}
        </Alert>
      </Snackbar>

      <Snackbar
        open={!!errorMessage}
        autoHideDuration={6000}
        onClose={() => setErrorMessage(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert severity="error" onClose={() => setErrorMessage(null)}>
          {errorMessage}
        </Alert>
      </Snackbar>
    </Container>
  );
};
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useForm, Controller } from 'react-hook-form';
import { useQuery } from '@apollo/client';
import {
  Container,
  Paper,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  CircularProgress,
} from '@mui/material';
import { useAuth } from '../contexts/AuthContext';
import { LIST_STORES } from '../graphql/queries';

interface LoginFormData {
  username: string;
  password: string;
  storeId: string;
}

export const Login: React.FC = () => {
  const navigate = useNavigate();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
  } = useForm<LoginFormData>();

  const { data: storesData, loading: storesLoading, error: storesError } = useQuery(LIST_STORES, {
    errorPolicy: 'all', // Allow partial data in case of authentication issues
  });

  const onSubmit = (data: LoginFormData) => {
    try {
      const success = login(data.username, data.password, data.storeId);
      if (success) {
        navigate('/');
      } else {
        setError('Invalid credentials');
      }
    } catch (err) {
      setError('Login failed. Please try again.');
    }
  };

  return (
    <Container maxWidth="sm">
      <Box sx={{ mt: 8 }}>
        <Paper elevation={3} sx={{ p: 4 }}>
          <Typography variant="h4" component="h1" gutterBottom align="center">
            Pet Store Login
          </Typography>
          <Typography variant="body2" color="text.secondary" align="center" sx={{ mb: 3 }}>
            Login as a customer to browse and purchase pets
          </Typography>
          
          {error && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {error}
            </Alert>
          )}

          <form onSubmit={handleSubmit(onSubmit)}>
            <TextField
              fullWidth
              label="Username"
              margin="normal"
              {...register('username', { required: 'Username is required' })}
              error={!!errors.username}
              helperText={errors.username?.message}
            />

            <TextField
              fullWidth
              label="Password"
              type="password"
              margin="normal"
              {...register('password', { required: 'Password is required' })}
              error={!!errors.password}
              helperText={errors.password?.message}
            />

            <FormControl fullWidth margin="normal" error={!!errors.storeId}>
              <InputLabel id="store-select-label">Select Store</InputLabel>
              <Controller
                name="storeId"
                control={control}
                rules={{ required: 'Please select a store' }}
                render={({ field }) => (
                  <Select
                    labelId="store-select-label"
                    label="Select Store"
                    {...field}
                    disabled={storesLoading}
                  >
                    {storesLoading ? (
                      <MenuItem disabled>
                        <CircularProgress size={20} sx={{ mr: 1 }} />
                        Loading stores...
                      </MenuItem>
                    ) : storesError ? (
                      <MenuItem disabled>
                        Error loading stores - using default store
                      </MenuItem>
                    ) : storesData?.listStores?.length ? (
                      storesData.listStores.map((store: any) => (
                        <MenuItem key={store.id} value={store.id}>
                          {store.name}
                        </MenuItem>
                      ))
                    ) : (
                      <MenuItem value="123e4567-e89b-12d3-a456-426614174000">
                        Pet Paradise Store (Default)
                      </MenuItem>
                    )}
                  </Select>
                )}
              />
              {errors.storeId && (
                <Typography variant="caption" color="error" sx={{ mt: 1 }}>
                  {errors.storeId.message}
                </Typography>
              )}
            </FormControl>

            <Button
              type="submit"
              fullWidth
              variant="contained"
              sx={{ mt: 3, mb: 2 }}
              size="large"
            >
              Login
            </Button>
          </form>

          <Typography variant="body2" color="text.secondary" align="center">
            Note: Use any username/password for customer access and select a store
          </Typography>
        </Paper>
      </Box>
    </Container>
  );
};
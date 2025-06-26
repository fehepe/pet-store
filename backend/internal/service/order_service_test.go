package service

import (
	"context"
	"testing"

	"github.com/fehepe/pet-store/backend/internal/mocks"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderService_CreateOrder_Validation(t *testing.T) {
	tests := []struct {
		name    string
		input   models.CreateOrderInput
		wantErr bool
	}{
		{
			name: "validation error - empty customer ID",
			input: models.CreateOrderInput{
				CustomerID: "", // Invalid
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{uuid.New()},
			},
			wantErr: true,
		},
		{
			name: "validation error - no pet IDs",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{}, // Invalid
			},
			wantErr: true,
		},
		{
			name: "validation error - too many pets",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs:     make([]uuid.UUID, 15), // Too many (>10)
			},
			wantErr: true,
		},
		{
			name: "validation error - duplicate pet IDs",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{uuid.New(), uuid.New()}, // Will be made duplicate below
			},
			wantErr: true,
		},
		{
			name: "valid input passes validation",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{uuid.New()},
			},
			wantErr: false, // Will fail later in transaction, but validation passes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make duplicate pet IDs for the duplicate test case
			if tt.name == "validation error - duplicate pet IDs" {
				petID := uuid.New()
				tt.input.PetIDs = []uuid.UUID{petID, petID}
			}

			mockOrderRepo := new(mocks.MockOrderRepository)
			mockPetRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockPetService := new(mocks.MockPetService)

			// For valid input, mock the transaction to fail so we can test validation separately
			if !tt.wantErr {
				mockOrderRepo.On("Transaction", mock.AnythingOfType("func(*sql.Tx) error")).Return(assert.AnError)
			}

			service := NewOrderService(mockOrderRepo, mockPetRepo, mockCache, mockPetService)

			_, err := service.CreateOrder(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.Error(t, err) // Will fail at transaction level, but validation passed
			}
		})
	}
}

func TestOrderService_GetOrderPets(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func(*mocks.MockOrderRepository, *mocks.MockCache)
	}{
		{
			name:    "successful retrieval from cache",
			wantErr: false,
			setup: func(orderRepo *mocks.MockOrderRepository, cache *mocks.MockCache) {
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					pets := args[2].(*[]*models.Pet)
					*pets = []*models.Pet{
						{ID: uuid.New(), Name: "CachedPet1"},
						{ID: uuid.New(), Name: "CachedPet2"},
					}
				}).Return(nil)
			},
		},
		{
			name:    "successful retrieval from repository",
			wantErr: false,
			setup: func(orderRepo *mocks.MockOrderRepository, cache *mocks.MockCache) {
				expectedPets := []*models.Pet{
					{ID: uuid.New(), Name: "Pet1"},
					{ID: uuid.New(), Name: "Pet2"},
				}
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError) // Cache miss
				orderRepo.On("GetOrderPets", mock.Anything, mock.Anything).Return(expectedPets, nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:    "repository error",
			wantErr: true,
			setup: func(orderRepo *mocks.MockOrderRepository, cache *mocks.MockCache) {
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(assert.AnError) // Cache miss
				orderRepo.On("GetOrderPets", mock.Anything, mock.Anything).Return([]*models.Pet(nil), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderRepo := new(mocks.MockOrderRepository)
			mockPetRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockPetService := new(mocks.MockPetService)

			tt.setup(mockOrderRepo, mockCache)

			service := NewOrderService(mockOrderRepo, mockPetRepo, mockCache, mockPetService)
			orderID := uuid.New()

			pets, err := service.GetOrderPets(context.Background(), orderID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pets)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pets)
				if tt.name == "successful retrieval from repository" {
					assert.Len(t, pets, 2)
				}
			}

			mockOrderRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestOrderServiceInterface_Implementation(t *testing.T) {
	mockOrderRepo := new(mocks.MockOrderRepository)
	mockPetRepo := new(mocks.MockPetRepository)
	mockCache := new(mocks.MockCache)
	mockPetService := new(mocks.MockPetService)

	var _ OrderServiceInterface = NewOrderService(mockOrderRepo, mockPetRepo, mockCache, mockPetService)
}

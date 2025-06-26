package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/fehepe/pet-store/backend/internal/mocks"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStoreService_CreateStore(t *testing.T) {
	tests := []struct {
		name    string
		input   models.CreateStoreInput
		wantErr bool
		setup   func(*mocks.MockStoreRepository, *mocks.MockCache)
	}{
		{
			name: "successful creation",
			input: models.CreateStoreInput{
				Name:    "Pet Paradise",
				OwnerID: "owner123",
			},
			wantErr: false,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache) {
				// No cache mocking needed for CreateStore - it calls repo directly
				// Mock GetByOwnerID to return not found (sql.ErrNoRows)
				repo.On("GetByOwnerID", mock.Anything, "owner123").Return(nil, sql.ErrNoRows)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Store")).Return(nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
			},
		},
		{
			name: "validation error - empty name",
			input: models.CreateStoreInput{
				Name:    "", // Invalid
				OwnerID: "owner123",
			},
			wantErr: true,
			setup:   func(*mocks.MockStoreRepository, *mocks.MockCache) {}, // No mocking needed for validation errors
		},
		{
			name: "validation error - empty owner ID",
			input: models.CreateStoreInput{
				Name:    "Pet Paradise",
				OwnerID: "", // Invalid
			},
			wantErr: true,
			setup: func(*mocks.MockStoreRepository, *mocks.MockCache) {}, // No mocking needed for validation errors
		},
		{
			name: "store already exists",
			input: models.CreateStoreInput{
				Name:    "Pet Paradise",
				OwnerID: "owner123",
			},
			wantErr: true,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache) {
				existingStore := &models.Store{
					ID:      uuid.New(),
					Name:    "Existing Store",
					OwnerID: "owner123",
				}
				// Mock direct repository call during CreateStore
				repo.On("GetByOwnerID", mock.Anything, "owner123").Return(existingStore, nil)
			},
		},
		{
			name: "repository creation error",
			input: models.CreateStoreInput{
				Name:    "Pet Paradise",
				OwnerID: "owner123",
			},
			wantErr: true,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache) {
				// Mock GetByOwnerID to return not found
				repo.On("GetByOwnerID", mock.Anything, "owner123").Return(nil, sql.ErrNoRows)
				// Mock Create to return error
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Store")).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockStoreRepository)
			mockCache := new(mocks.MockCache)

			tt.setup(mockRepo, mockCache)

			service := NewStoreService(mockRepo, mockCache)

			store, err := service.CreateStore(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, store)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, store)
				assert.Equal(t, tt.input.Name, store.Name)
				assert.Equal(t, tt.input.OwnerID, store.OwnerID)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestStoreService_GetStoreByOwnerID(t *testing.T) {
	tests := []struct {
		name    string
		ownerID string
		wantErr bool
		setup   func(*mocks.MockStoreRepository, *mocks.MockCache)
	}{
		{
			name:    "successful retrieval from cache",
			ownerID: "owner123",
			wantErr: false,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache) {
				// Mock cache hit with actual store data
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					store := args[2].(*models.Store)
					store.ID = uuid.New()
					store.Name = "Cached Store"
					store.OwnerID = "owner123"
				}).Return(nil)
			},
		},
		{
			name:    "successful retrieval from repository",
			ownerID: "owner123",
			wantErr: false,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache) {
				expectedStore := &models.Store{
					ID:      uuid.New(),
					Name:    "Pet Paradise",
					OwnerID: "owner123",
				}
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows) // Cache miss
				repo.On("GetByOwnerID", mock.Anything, "owner123").Return(expectedStore, nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:    "store not found",
			ownerID: "newowner",
			wantErr: true,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache) {
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows) // Cache miss
				repo.On("GetByOwnerID", mock.Anything, "newowner").Return(nil, sql.ErrNoRows)
			},
		},
		{
			name:    "validation error - empty owner ID",
			ownerID: "",
			wantErr: true,
			setup:   func(*mocks.MockStoreRepository, *mocks.MockCache) {}, // No mocking needed for validation errors
		},
		{
			name:    "validation error - whitespace only owner ID",
			ownerID: "   ",
			wantErr: true,
			setup:   func(*mocks.MockStoreRepository, *mocks.MockCache) {}, // No mocking needed for validation errors
		},
		{
			name:    "repository error",
			ownerID: "owner123",
			wantErr: true,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache) {
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows) // Cache miss
				repo.On("GetByOwnerID", mock.Anything, "owner123").Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockStoreRepository)
			mockCache := new(mocks.MockCache)

			tt.setup(mockRepo, mockCache)

			service := NewStoreService(mockRepo, mockCache)

			store, err := service.GetStoreByOwnerID(context.Background(), tt.ownerID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, store)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, store)
				assert.Equal(t, tt.ownerID, store.OwnerID)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestStoreService_GetStoreByID(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func(*mocks.MockStoreRepository, *mocks.MockCache, uuid.UUID)
	}{
		{
			name:    "successful retrieval from cache",
			wantErr: false,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache, storeID uuid.UUID) {
				// Mock cache hit with actual store data
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
					store := args[2].(*models.Store)
					store.ID = storeID
					store.Name = "Cached Store"
					store.OwnerID = "owner123"
				}).Return(nil)
			},
		},
		{
			name:    "successful retrieval from repository",
			wantErr: false,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache, storeID uuid.UUID) {
				expectedStore := &models.Store{
					ID:      storeID,
					Name:    "Pet Paradise",
					OwnerID: "owner123",
				}
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows) // Cache miss
				repo.On("GetByID", mock.Anything, storeID).Return(expectedStore, nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:    "store not found",
			wantErr: true,
			setup: func(repo *mocks.MockStoreRepository, cache *mocks.MockCache, storeID uuid.UUID) {
				cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows) // Cache miss
				repo.On("GetByID", mock.Anything, storeID).Return(nil, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockStoreRepository)
			mockCache := new(mocks.MockCache)
			storeID := uuid.New()

			tt.setup(mockRepo, mockCache, storeID)

			service := NewStoreService(mockRepo, mockCache)

			store, err := service.GetStoreByID(context.Background(), storeID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, store)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, store)
				assert.Equal(t, storeID, store.ID)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestStoreServiceInterface_Implementation(t *testing.T) {
	// Test that StoreService implements StoreServiceInterface
	mockRepo := new(mocks.MockStoreRepository)
	mockCache := new(mocks.MockCache)

	var _ StoreServiceInterface = NewStoreService(mockRepo, mockCache)
}
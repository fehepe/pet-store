package mocks

import (
	"context"

	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockStoreRepository is a mock implementation of StoreRepositoryInterface
type MockStoreRepository struct {
	mock.Mock
}

func (m *MockStoreRepository) Create(ctx context.Context, store *models.Store) error {
	args := m.Called(ctx, store)
	return args.Error(0)
}

func (m *MockStoreRepository) GetByID(ctx context.Context, storeID uuid.UUID) (*models.Store, error) {
	args := m.Called(ctx, storeID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Store), args.Error(1)
}

func (m *MockStoreRepository) GetByOwnerID(ctx context.Context, ownerID string) (*models.Store, error) {
	args := m.Called(ctx, ownerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Store), args.Error(1)
}


package mocks

import (
	"context"
	"database/sql"

	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockPetRepository is a mock implementation of PetRepositoryInterface
type MockPetRepository struct {
	mock.Mock
}

func (m *MockPetRepository) Create(ctx context.Context, pet *models.Pet) error {
	args := m.Called(ctx, pet)
	return args.Error(0)
}

func (m *MockPetRepository) GetByID(ctx context.Context, petID uuid.UUID) (*models.Pet, error) {
	args := m.Called(ctx, petID)
	return args.Get(0).(*models.Pet), args.Error(1)
}

func (m *MockPetRepository) List(ctx context.Context, filter models.PetFilter) ([]*models.Pet, int, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*models.Pet), args.Int(1), args.Error(2)
}

func (m *MockPetRepository) Delete(ctx context.Context, petID uuid.UUID) error {
	args := m.Called(ctx, petID)
	return args.Error(0)
}

func (m *MockPetRepository) MarkAsSold(ctx context.Context, tx *sql.Tx, petID uuid.UUID) error {
	args := m.Called(ctx, tx, petID)
	return args.Error(0)
}

func (m *MockPetRepository) Transaction(fn func(*sql.Tx) error) error {
	args := m.Called(fn)
	return args.Error(0)
}
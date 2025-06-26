package mocks

import (
	"context"

	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockPetService is a mock implementation of PetServiceInterface for order tests
type MockPetService struct {
	mock.Mock
}

func (m *MockPetService) CreatePet(ctx context.Context, input models.CreatePetInput) (*models.Pet, error) {
	args := m.Called(ctx, input)
	return args.Get(0).(*models.Pet), args.Error(1)
}

func (m *MockPetService) GetPetByID(ctx context.Context, petID uuid.UUID) (*models.Pet, error) {
	args := m.Called(ctx, petID)
	return args.Get(0).(*models.Pet), args.Error(1)
}

func (m *MockPetService) ListPets(ctx context.Context, filter models.PetFilter) ([]*models.Pet, int, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*models.Pet), args.Int(1), args.Error(2)
}

func (m *MockPetService) DeletePetByID(ctx context.Context, petID uuid.UUID) error {
	args := m.Called(ctx, petID)
	return args.Error(0)
}

func (m *MockPetService) MarkPetAsSold(ctx context.Context, petID uuid.UUID) error {
	args := m.Called(ctx, petID)
	return args.Error(0)
}

func (m *MockPetService) DecryptBreederEmail(encryptedEmail string) (string, error) {
	args := m.Called(encryptedEmail)
	return args.String(0), args.Error(1)
}
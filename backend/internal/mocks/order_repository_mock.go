package mocks

import (
	"context"
	"database/sql"

	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository is a mock implementation of OrderRepositoryInterface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, order *models.Order) error {
	args := m.Called(ctx, tx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, order *models.Order) error {
	args := m.Called(ctx, tx, order)
	return args.Error(0)
}

func (m *MockOrderRepository) CreateItem(ctx context.Context, tx *sql.Tx, item *models.OrderItem) error {
	args := m.Called(ctx, tx, item)
	return args.Error(0)
}

func (m *MockOrderRepository) GetOrderPets(ctx context.Context, orderID uuid.UUID) ([]*models.Pet, error) {
	args := m.Called(ctx, orderID)
	return args.Get(0).([]*models.Pet), args.Error(1)
}

func (m *MockOrderRepository) Transaction(fn func(*sql.Tx) error) error {
	args := m.Called(fn)
	return args.Error(0)
}


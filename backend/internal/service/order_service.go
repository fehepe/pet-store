package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fehepe/pet-store/backend/internal/cache"
	apperrors "github.com/fehepe/pet-store/backend/internal/errors"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/fehepe/pet-store/backend/internal/repository"
	"github.com/fehepe/pet-store/backend/internal/validation"
	"github.com/google/uuid"
)

// OrderServiceInterface defines the interface for order operations
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, input models.CreateOrderInput) (*models.Order, error)
	GetOrderPets(ctx context.Context, orderID uuid.UUID) ([]*models.Pet, error)
}

// OrderService implements OrderServiceInterface with improved error handling and validation
type OrderService struct {
	repo       repository.OrderRepositoryInterface
	petRepo    repository.PetRepositoryInterface
	cache      cache.CacheInterface
	petService PetServiceInterface
}

// NewOrderService creates a new order service
func NewOrderService(
	repo repository.OrderRepositoryInterface,
	petRepo repository.PetRepositoryInterface,
	cache cache.CacheInterface,
	petService PetServiceInterface,
) *OrderService {
	return &OrderService{
		repo:       repo,
		petRepo:    petRepo,
		cache:      cache,
		petService: petService,
	}
}

// CreateOrder creates a new order with proper validation and error handling
func (s *OrderService) CreateOrder(ctx context.Context, input models.CreateOrderInput) (*models.Order, error) {
	if err := validation.ValidateCreateOrderInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	input.CustomerID = validation.SanitizeString(input.CustomerID)

	var order *models.Order
	var orderItems []*models.OrderItem
	unavailablePets := []string{}

	err := s.repo.Transaction(func(tx *sql.Tx) error {
		order = &models.Order{
			ID:         uuid.New(),
			CustomerID: input.CustomerID,
			StoreID:    input.StoreID,
			TotalPets:  len(input.PetIDs),
			CreatedAt:  time.Now(),
		}

		if err := s.repo.CreateWithTx(ctx, tx, order); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		for _, petID := range input.PetIDs {
			var petName string
			checkQuery := `SELECT name FROM pets WHERE id = $1 AND store_id = $2 AND status = $3 FOR UPDATE`
			err := tx.QueryRowContext(ctx, checkQuery, petID, input.StoreID, models.PetStatusAvailable).Scan(&petName)
			if err == sql.ErrNoRows {
				unavailablePets = append(unavailablePets, petID.String())
				continue
			} else if err != nil {
				return fmt.Errorf("failed to check pet availability: %w", err)
			}

			if err := s.petRepo.MarkAsSold(ctx, tx, petID); err != nil {
				unavailablePets = append(unavailablePets, petName)
				continue
			}

			orderItem := &models.OrderItem{
				ID:          uuid.New(),
				OrderID:     order.ID,
				PetID:       petID,
				PurchasedAt: time.Now(),
			}

			if err := s.repo.CreateItem(ctx, tx, orderItem); err != nil {
				return fmt.Errorf("failed to create order item: %w", err)
			}

			orderItems = append(orderItems, orderItem)
		}

		if len(orderItems) == 0 {
			return apperrors.NewBusinessRuleError("no pets were available for purchase")
		}

		if len(orderItems) != len(input.PetIDs) {
			order.TotalPets = len(orderItems)
			if err := s.repo.UpdateWithTx(ctx, tx, order); err != nil {
				return fmt.Errorf("failed to update order: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, item := range orderItems {
		_ = s.cache.Delete(ctx, cache.PetCacheKey(input.StoreID.String(), item.PetID.String()))
	}
	_ = s.cache.InvalidatePattern(ctx, fmt.Sprintf("pets:list:%s:*", input.StoreID))

	if len(unavailablePets) > 0 {
		return order, fmt.Errorf("the following pets are no longer available: %v", unavailablePets)
	}

	return order, nil
}

// GetOrderPets retrieves pets for a specific order
func (s *OrderService) GetOrderPets(ctx context.Context, orderID uuid.UUID) ([]*models.Pet, error) {
	cacheKey := fmt.Sprintf("order:pets:%s", orderID.String())
	var pets []*models.Pet
	if err := s.cache.Get(ctx, cacheKey, &pets); err == nil {
		return pets, nil
	}

	pets, err := s.repo.GetOrderPets(ctx, orderID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, pets, 10*time.Minute)

	return pets, nil
}

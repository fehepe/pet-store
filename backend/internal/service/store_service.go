package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/fehepe/pet-store/backend/internal/cache"
	apperrors "github.com/fehepe/pet-store/backend/internal/errors"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/fehepe/pet-store/backend/internal/repository"
	"github.com/fehepe/pet-store/backend/internal/validation"
	"github.com/google/uuid"
)

// StoreServiceInterface defines the interface for store operations
type StoreServiceInterface interface {
	CreateStore(ctx context.Context, input models.CreateStoreInput) (*models.Store, error)
	GetStoreByOwnerID(ctx context.Context, ownerID string) (*models.Store, error)
}

// StoreService implements StoreServiceInterface with improved error handling and validation
type StoreService struct {
	repo  repository.StoreRepositoryInterface
	cache cache.CacheInterface
}

// NewStoreService creates a new store service
func NewStoreService(
	repo repository.StoreRepositoryInterface,
	cache cache.CacheInterface,
) *StoreService {
	return &StoreService{
		repo:  repo,
		cache: cache,
	}
}

// CreateStore creates a new store with proper validation and error handling
func (s *StoreService) CreateStore(ctx context.Context, input models.CreateStoreInput) (*models.Store, error) {
	if err := validation.ValidateCreateStoreInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	input.Name = validation.SanitizeString(input.Name)
	input.OwnerID = validation.SanitizeString(input.OwnerID)

	// Check if store already exists for this owner (direct repository call to avoid recursion)
	existingStore, err := s.repo.GetByOwnerID(ctx, input.OwnerID)
	if err == nil && existingStore != nil {
		return nil, apperrors.ConflictError{
			Resource: "store",
			Message:  "store already exists for this owner",
		}
	}

	store := &models.Store{
		ID:        uuid.New(),
		Name:      input.Name,
		OwnerID:   input.OwnerID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.repo.Create(ctx, store)
	if err != nil {
		return nil, fmt.Errorf("failed to create store: %w", err)
	}

	cacheKey := cache.StoreCacheKey(store.ID.String())
	_ = s.cache.Set(ctx, cacheKey, store, 10*time.Minute)

	ownerCacheKey := fmt.Sprintf("store:owner:%s", store.OwnerID)
	_ = s.cache.Set(ctx, ownerCacheKey, store, 10*time.Minute)

	return store, nil
}

// GetStoreByOwnerID retrieves a store by owner ID with caching
func (s *StoreService) GetStoreByOwnerID(ctx context.Context, ownerID string) (*models.Store, error) {
	if strings.TrimSpace(ownerID) == "" {
		return nil, apperrors.NewValidationError("ownerID", "owner ID cannot be empty")
	}

	cacheKey := fmt.Sprintf("store:owner:%s", ownerID)
	var store models.Store
	if err := s.cache.Get(ctx, cacheKey, &store); err == nil {
		return &store, nil
	}

	storePtr, err := s.repo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, storePtr, 10*time.Minute)

	return storePtr, nil
}


func (s *StoreService) ListAllStores(ctx context.Context) ([]*models.Store, error) {
	stores, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return stores, nil
}


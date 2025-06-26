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
	"github.com/fehepe/pet-store/backend/pkg/encryption"
	"github.com/google/uuid"
)

// PetServiceInterface defines the interface for pet operations
type PetServiceInterface interface {
	CreatePet(ctx context.Context, input models.CreatePetInput) (*models.Pet, error)
	GetPetByID(ctx context.Context, petID uuid.UUID) (*models.Pet, error)
	ListPets(ctx context.Context, filter models.PetFilter) ([]*models.Pet, int, error)
	DeletePetByID(ctx context.Context, petID uuid.UUID) error
	MarkPetAsSold(ctx context.Context, petID uuid.UUID) error
	DecryptBreederEmail(encryptedEmail string) (string, error)
}

// PetService implements PetServiceInterface with improved error handling and validation
type PetService struct {
	repo      repository.PetRepositoryInterface
	cache     cache.CacheInterface
	encryptor encryption.EncryptorInterface
}

// NewPetService creates a new pet service
func NewPetService(
	repo repository.PetRepositoryInterface,
	cache cache.CacheInterface,
	encryptor encryption.EncryptorInterface,
) *PetService {
	return &PetService{
		repo:      repo,
		cache:     cache,
		encryptor: encryptor,
	}
}

// CreatePet creates a new pet with proper validation and error handling
func (s *PetService) CreatePet(ctx context.Context, input models.CreatePetInput) (*models.Pet, error) {
	if err := validation.ValidateCreatePetInput(input); err != nil {
		return nil, fmt.Errorf("invalid input: %w", err)
	}

	input.Name = validation.SanitizeString(input.Name)
	input.BreederName = validation.SanitizeString(input.BreederName)
	input.BreederEmail = validation.SanitizeString(input.BreederEmail)
	if input.Description != nil {
		desc := validation.SanitizeString(*input.Description)
		input.Description = &desc
	}

	encryptedEmail, err := s.encryptor.Encrypt(input.BreederEmail)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt breeder email: %w", err)
	}

	pet := &models.Pet{
		ID:                    uuid.New(),
		StoreID:               input.StoreID,
		Name:                  input.Name,
		Species:               input.Species,
		Age:                   input.Age,
		PictureURL:            input.PictureURL,
		Description:           input.Description,
		BreederName:           input.BreederName,
		BreederEmailEncrypted: encryptedEmail,
		Status:                models.PetStatusAvailable,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	err = s.repo.Create(ctx, pet)
	if err != nil {
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}

	cacheKey := cache.PetCacheKey(pet.StoreID.String(), pet.ID.String())
	_ = s.cache.Set(ctx, cacheKey, pet, 5*time.Minute)
	_ = s.cache.InvalidatePattern(ctx, fmt.Sprintf("pets:list:%s:*", pet.StoreID))

	return pet, nil
}

// GetPetByID retrieves a pet by ID with caching
func (s *PetService) GetPetByID(ctx context.Context, petID uuid.UUID) (*models.Pet, error) {
	pet, err := s.repo.GetByID(ctx, petID)
	if err != nil {
		return nil, err
	}

	cacheKey := cache.PetCacheKey(pet.StoreID.String(), pet.ID.String())
	_ = s.cache.Set(ctx, cacheKey, pet, 5*time.Minute)

	return pet, nil
}

// ListPets retrieves pets with filtering and pagination
func (s *PetService) ListPets(ctx context.Context, filter models.PetFilter) ([]*models.Pet, int, error) {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}

	pets, totalCount, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return pets, totalCount, nil
}

// DeletePetByID deletes a pet by ID
func (s *PetService) DeletePetByID(ctx context.Context, petID uuid.UUID) error {
	// First check if pet exists and get its store ID for cache invalidation
	pet, err := s.GetPetByID(ctx, petID)
	if err != nil {
		return err // Already returns proper error type from GetPetByID
	}

	// Check if pet can be deleted (not sold)
	if pet.Status == models.PetStatusSold {
		return apperrors.ConflictError{
			Resource: "pet",
			Message:  "cannot delete a sold pet",
		}
	}

	err = s.repo.Delete(ctx, petID)
	if err != nil {
		return err
	}
	cacheKey := cache.PetCacheKey(pet.StoreID.String(), petID.String())
	_ = s.cache.Delete(ctx, cacheKey)
	_ = s.cache.InvalidatePattern(ctx, fmt.Sprintf("pets:list:%s:*", pet.StoreID))

	return nil
}

// MarkPetAsSold marks a pet as sold (creates its own transaction)
func (s *PetService) MarkPetAsSold(ctx context.Context, petID uuid.UUID) error {
	// For standalone usage, create a transaction
	return s.repo.Transaction(func(tx *sql.Tx) error {
		return s.repo.MarkAsSold(ctx, tx, petID)
	})
}


// DecryptBreederEmail decrypts the breeder email
func (s *PetService) DecryptBreederEmail(encryptedEmail string) (string, error) {
	if encryptedEmail == "" {
		return "", fmt.Errorf("encrypted email is empty")
	}

	decrypted, err := s.encryptor.Decrypt(encryptedEmail)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt breeder email: %w", err)
	}

	return decrypted, nil
}

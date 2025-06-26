package validation

import (
	"regexp"
	"strings"

	apperrors "github.com/fehepe/pet-store/backend/internal/errors"
	"github.com/fehepe/pet-store/backend/internal/models"
)

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// ValidateCreatePetInput validates the input for creating a pet
func ValidateCreatePetInput(input models.CreatePetInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return apperrors.NewValidationError("name", "pet name is required and cannot be empty")
	}

	if len(strings.TrimSpace(input.Name)) > 100 {
		return apperrors.NewValidationError("name", "pet name cannot exceed 100 characters")
	}

	if input.Age < 0 {
		return apperrors.NewValidationError("age", "pet age cannot be negative")
	}

	if input.Age > 50 {
		return apperrors.NewValidationError("age", "pet age cannot exceed 50 years")
	}

	if !IsValidSpecies(string(input.Species)) {
		return apperrors.NewValidationError("species", "species must be Cat, Dog, or Frog")
	}

	if strings.TrimSpace(input.BreederName) == "" {
		return apperrors.NewValidationError("breederName", "breeder name is required and cannot be empty")
	}

	if len(strings.TrimSpace(input.BreederName)) > 100 {
		return apperrors.NewValidationError("breederName", "breeder name cannot exceed 100 characters")
	}

	if !IsValidEmail(input.BreederEmail) {
		return apperrors.NewValidationError("breederEmail", "breeder email must be a valid email address")
	}

	if input.Description != nil && len(*input.Description) > 1000 {
		return apperrors.NewValidationError("description", "description cannot exceed 1000 characters")
	}

	if input.PictureURL != nil && len(*input.PictureURL) > 500 {
		return apperrors.NewValidationError("pictureURL", "picture URL cannot exceed 500 characters")
	}

	return nil
}

// IsValidEmail checks if the email format is valid
func IsValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if email == "" {
		return false
	}
	if len(email) > 254 { // RFC 5321 limit
		return false
	}
	return emailRegex.MatchString(email)
}

// IsValidSpecies checks if the species is valid
func IsValidSpecies(species string) bool {
	validSpecies := map[string]bool{
		"Cat":  true,
		"Dog":  true,
		"Frog": true,
	}
	return validSpecies[species]
}

// SanitizeString removes dangerous characters and trims whitespace
func SanitizeString(s string) string {
	// Remove null bytes and trim whitespace
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.TrimSpace(s)
	return s
}

// ValidateCreateStoreInput validates the input for creating a store
func ValidateCreateStoreInput(input models.CreateStoreInput) error {
	if strings.TrimSpace(input.Name) == "" {
		return apperrors.NewValidationError("name", "store name is required and cannot be empty")
	}

	if len(strings.TrimSpace(input.Name)) > 100 {
		return apperrors.NewValidationError("name", "store name cannot exceed 100 characters")
	}

	if strings.TrimSpace(input.OwnerID) == "" {
		return apperrors.NewValidationError("ownerID", "owner ID is required and cannot be empty")
	}

	if len(strings.TrimSpace(input.OwnerID)) > 50 {
		return apperrors.NewValidationError("ownerID", "owner ID cannot exceed 50 characters")
	}

	return nil
}

// ValidateCreateOrderInput validates the input for creating an order
func ValidateCreateOrderInput(input models.CreateOrderInput) error {
	if strings.TrimSpace(input.CustomerID) == "" {
		return apperrors.NewValidationError("customerID", "customer ID is required and cannot be empty")
	}

	if len(strings.TrimSpace(input.CustomerID)) > 50 {
		return apperrors.NewValidationError("customerID", "customer ID cannot exceed 50 characters")
	}

	if len(input.PetIDs) == 0 {
		return apperrors.NewValidationError("petIDs", "at least one pet ID is required")
	}

	if len(input.PetIDs) > 10 {
		return apperrors.NewValidationError("petIDs", "cannot purchase more than 10 pets in a single order")
	}

	// Check for duplicate pet IDs
	petIDMap := make(map[string]bool)
	for _, petID := range input.PetIDs {
		petIDStr := petID.String()
		if petIDStr == "00000000-0000-0000-0000-000000000000" {
			return apperrors.NewValidationError("petIDs", "pet ID cannot be empty")
		}
		if petIDMap[petIDStr] {
			return apperrors.NewValidationError("petIDs", "duplicate pet IDs are not allowed")
		}
		petIDMap[petIDStr] = true
	}

	return nil
}

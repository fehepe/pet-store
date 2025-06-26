package errors

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNotFoundError(t *testing.T) {
	err := NotFoundError{
		Resource: "pet",
		ID:       "123",
	}
	expected := "pet with ID 123 not found"
	assert.Equal(t, expected, err.Error())
}

func TestValidationError(t *testing.T) {
	err := ValidationError{
		Field:   "name",
		Message: "cannot be empty",
	}
	expected := "validation error for name: cannot be empty"
	assert.Equal(t, expected, err.Error())
}

func TestConflictError(t *testing.T) {
	err := ConflictError{
		Resource: "pet",
		Message:  "already sold",
	}
	expected := "conflict with pet: already sold"
	assert.Equal(t, expected, err.Error())
}

func TestPetNotFoundError(t *testing.T) {
	petID := uuid.New()
	err := PetNotFoundError{
		PetID: petID,
	}
	expected := "pet with ID " + petID.String() + " not found"
	assert.Equal(t, expected, err.Error())
}

func TestStoreNotFoundError(t *testing.T) {
	storeID := uuid.New()
	err := StoreNotFoundError{
		StoreID: storeID,
	}
	expected := "store with ID " + storeID.String() + " not found"
	assert.Equal(t, expected, err.Error())
}

func TestNewPetNotFound(t *testing.T) {
	petID := uuid.New()
	err := NewPetNotFound(petID)
	assert.IsType(t, PetNotFoundError{}, err)
	assert.Contains(t, err.Error(), petID.String())
}

func TestNewStoreNotFound(t *testing.T) {
	storeID := uuid.New()
	err := NewStoreNotFound(storeID)
	assert.IsType(t, StoreNotFoundError{}, err)
	assert.Contains(t, err.Error(), storeID.String())
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("name", "cannot be empty")
	assert.IsType(t, ValidationError{}, err)
	assert.Contains(t, err.Error(), "name")
	assert.Contains(t, err.Error(), "cannot be empty")
}


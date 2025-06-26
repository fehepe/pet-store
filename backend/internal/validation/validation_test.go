package validation

import (
	"testing"

	apperrors "github.com/fehepe/pet-store/backend/internal/errors"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreatePetInput(t *testing.T) {
	tests := []struct {
		name      string
		input     models.CreatePetInput
		wantError bool
		errorType interface{}
	}{
		{
			name: "valid input",
			input: models.CreatePetInput{
				StoreID:      uuid.New(),
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantError: false,
		},
		{
			name: "empty name",
			input: models.CreatePetInput{
				Name:         "",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "name too long",
			input: models.CreatePetInput{
				Name:         "This is a very long name that exceeds the maximum allowed length of one hundred characters for any pet",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "negative age",
			input: models.CreatePetInput{
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          -1,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "age too high",
			input: models.CreatePetInput{
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          51,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "invalid species",
			input: models.CreatePetInput{
				Name:         "Fluffy",
				Species:      models.PetSpecies("Bird"), // Invalid species
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "empty breeder name",
			input: models.CreatePetInput{
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "",
				BreederEmail: "john@example.com",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "invalid email",
			input: models.CreatePetInput{
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "invalid-email",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreatePetInput(tt.input)

			if tt.wantError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{"valid email", "test@example.com", true},
		{"valid email with subdomain", "test@mail.example.com", true},
		{"valid email with numbers", "test123@example.com", true},
		{"valid email with special chars", "test.user+tag@example.com", true},
		{"empty email", "", false},
		{"email without @", "testexample.com", false},
		{"email without domain", "test@", false},
		{"email without TLD", "test@example", false},
		{"email with spaces", "test @example.com", false},
		{"email too long", "a" + string(make([]byte, 250)) + "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestIsValidSpecies(t *testing.T) {
	tests := []struct {
		name    string
		species string
		want    bool
	}{
		{"Cat", "Cat", true},
		{"Dog", "Dog", true},
		{"Frog", "Frog", true},
		{"cat lowercase", "cat", false},
		{"Bird", "Bird", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidSpecies(tt.species)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestSanitizeString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		{"normal string", "hello world", "hello world"},
		{"string with null bytes", "hello\x00world", "helloworld"},
		{"string with leading/trailing spaces", "  hello world  ", "hello world"},
		{"empty string", "", ""},
		{"only spaces", "   ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeString(tt.input)
			assert.Equal(t, tt.output, result)
		})
	}
}

func TestValidateCreateStoreInput(t *testing.T) {
	tests := []struct {
		name      string
		input     models.CreateStoreInput
		wantError bool
		errorType interface{}
	}{
		{
			name: "valid input",
			input: models.CreateStoreInput{
				Name:    "Pet Paradise",
				OwnerID: "owner123",
			},
			wantError: false,
		},
		{
			name: "empty name",
			input: models.CreateStoreInput{
				Name:    "",
				OwnerID: "owner123",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "name too long",
			input: models.CreateStoreInput{
				Name:    "This is a very long store name that exceeds one hundred characters in length and should fail validation",
				OwnerID: "owner123",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "empty owner ID",
			input: models.CreateStoreInput{
				Name:    "Pet Paradise",
				OwnerID: "",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "owner ID too long",
			input: models.CreateStoreInput{
				Name:    "Pet Paradise",
				OwnerID: "this_is_a_very_long_owner_id_that_exceeds_fifty_chars",
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateStoreInput(tt.input)

			if tt.wantError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateCreateOrderInput(t *testing.T) {
	tests := []struct {
		name      string
		input     models.CreateOrderInput
		wantError bool
		errorType interface{}
	}{
		{
			name: "valid input",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{uuid.New(), uuid.New()},
			},
			wantError: false,
		},
		{
			name: "empty customer ID",
			input: models.CreateOrderInput{
				CustomerID: "",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{uuid.New()},
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "customer ID too long",
			input: models.CreateOrderInput{
				CustomerID: "this_is_a_very_long_customer_id_that_exceeds_fifty_characters_limit",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{uuid.New()},
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "empty pet IDs",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{},
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "too many pet IDs",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs: []uuid.UUID{
					uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(),
					uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(),
					uuid.New(), // 11 pets - should fail
				},
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "duplicate pet IDs",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs: func() []uuid.UUID {
					petID := uuid.New()
					return []uuid.UUID{petID, petID} // Duplicate
				}(),
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
		{
			name: "empty pet ID in list",
			input: models.CreateOrderInput{
				CustomerID: "customer123",
				StoreID:    uuid.New(),
				PetIDs:     []uuid.UUID{uuid.New(), uuid.Nil},
			},
			wantError: true,
			errorType: apperrors.ValidationError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCreateOrderInput(tt.input)

			if tt.wantError {
				assert.Error(t, err)
				assert.IsType(t, tt.errorType, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Helper function to create int32 pointer
func int32Ptr(i int32) *int32 {
	return &i
}

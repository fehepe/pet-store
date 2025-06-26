package service

import (
	"context"
	"testing"

	"github.com/fehepe/pet-store/backend/internal/mocks"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPetService_CreatePet_WithMocks(t *testing.T) {
	tests := []struct {
		name    string
		input   models.CreatePetInput
		wantErr bool
		setup   func(*mocks.MockPetRepository, *mocks.MockCache, *mocks.MockEncryptor)
	}{
		{
			name: "successful creation",
			input: models.CreatePetInput{
				StoreID:      uuid.New(),
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantErr: false,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, encryptor *mocks.MockEncryptor) {
				encryptor.On("Encrypt", "john@example.com").Return("encrypted_email", nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Pet")).Return(nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				cache.On("InvalidatePattern", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name: "validation error - empty name",
			input: models.CreatePetInput{
				StoreID:      uuid.New(),
				Name:         "", // Invalid
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantErr: true,
			setup:   func(*mocks.MockPetRepository, *mocks.MockCache, *mocks.MockEncryptor) {}, // No mocking needed for validation errors
		},
		{
			name: "validation error - invalid email",
			input: models.CreatePetInput{
				StoreID:      uuid.New(),
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "invalid-email", // Invalid
			},
			wantErr: true,
			setup:   func(*mocks.MockPetRepository, *mocks.MockCache, *mocks.MockEncryptor) {}, // No mocking needed for validation errors
		},
		{
			name: "encryption error",
			input: models.CreatePetInput{
				StoreID:      uuid.New(),
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantErr: true,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, encryptor *mocks.MockEncryptor) {
				encryptor.On("Encrypt", "john@example.com").Return("", assert.AnError)
			},
		},
		{
			name: "repository creation error",
			input: models.CreatePetInput{
				StoreID:      uuid.New(),
				Name:         "Fluffy",
				Species:      models.PetSpeciesCat,
				Age:          3,
				BreederName:  "John Doe",
				BreederEmail: "john@example.com",
			},
			wantErr: true,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, encryptor *mocks.MockEncryptor) {
				encryptor.On("Encrypt", "john@example.com").Return("encrypted_email", nil)
				repo.On("Create", mock.Anything, mock.AnythingOfType("*models.Pet")).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockEncryptor := new(mocks.MockEncryptor)

			tt.setup(mockRepo, mockCache, mockEncryptor)

			service := NewPetService(mockRepo, mockCache, mockEncryptor)

			pet, err := service.CreatePet(context.Background(), tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pet)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pet)
				assert.Equal(t, tt.input.Name, pet.Name)
				assert.Equal(t, tt.input.Species, pet.Species)
				assert.Equal(t, models.PetStatusAvailable, pet.Status)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
			mockEncryptor.AssertExpectations(t)
		})
	}
}

func TestPetService_GetPetByID(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func(*mocks.MockPetRepository, *mocks.MockCache, uuid.UUID)
	}{
		{
			name:    "successful retrieval",
			wantErr: false,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, petID uuid.UUID) {
				expectedPet := &models.Pet{
					ID:      petID,
					StoreID: uuid.New(),
					Name:    "Fluffy",
					Species: models.PetSpeciesCat,
				}
				repo.On("GetByID", mock.Anything, petID).Return(expectedPet, nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:    "pet not found",
			wantErr: true,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, petID uuid.UUID) {
				repo.On("GetByID", mock.Anything, petID).Return((*models.Pet)(nil), assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			petID := uuid.New()
			mockRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockEncryptor := new(mocks.MockEncryptor)

			tt.setup(mockRepo, mockCache, petID)

			service := NewPetService(mockRepo, mockCache, mockEncryptor)

			pet, err := service.GetPetByID(context.Background(), petID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pet)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pet)
				assert.Equal(t, petID, pet.ID)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestPetService_ListPets(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func(*mocks.MockPetRepository)
	}{
		{
			name:    "successful listing",
			wantErr: false,
			setup: func(repo *mocks.MockPetRepository) {
				expectedPets := []*models.Pet{
					{ID: uuid.New(), Name: "Pet1"},
					{ID: uuid.New(), Name: "Pet2"},
				}
				repo.On("List", mock.Anything, mock.AnythingOfType("models.PetFilter")).Return(expectedPets, 2, nil)
			},
		},
		{
			name:    "repository error",
			wantErr: true,
			setup: func(repo *mocks.MockPetRepository) {
				repo.On("List", mock.Anything, mock.AnythingOfType("models.PetFilter")).Return([]*models.Pet(nil), 0, assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := models.PetFilter{
				Limit:  10,
				Offset: 0,
			}

			mockRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockEncryptor := new(mocks.MockEncryptor)

			tt.setup(mockRepo)

			service := NewPetService(mockRepo, mockCache, mockEncryptor)

			pets, total, err := service.ListPets(context.Background(), filter)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, pets)
				assert.Equal(t, 0, total)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, pets)
				assert.Equal(t, 2, len(pets))
				assert.Equal(t, 2, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPetService_DeletePetByID(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func(*mocks.MockPetRepository, *mocks.MockCache, *mocks.MockEncryptor)
	}{
		{
			name:    "successful deletion",
			wantErr: false,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, encryptor *mocks.MockEncryptor) {
				pet := &models.Pet{
					ID:      uuid.New(),
					StoreID: uuid.New(),
					Name:    "Fluffy",
					Status:  models.PetStatusAvailable,
				}
				repo.On("GetByID", mock.Anything, mock.Anything).Return(pet, nil)
				repo.On("Delete", mock.Anything, mock.Anything).Return(nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				cache.On("Delete", mock.Anything, mock.Anything).Return(nil)
				cache.On("InvalidatePattern", mock.Anything, mock.Anything).Return(nil)
			},
		},
		{
			name:    "pet not found",
			wantErr: true,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, encryptor *mocks.MockEncryptor) {
				repo.On("GetByID", mock.Anything, mock.Anything).Return((*models.Pet)(nil), assert.AnError)
			},
		},
		{
			name:    "deletion error",
			wantErr: true,
			setup: func(repo *mocks.MockPetRepository, cache *mocks.MockCache, encryptor *mocks.MockEncryptor) {
				pet := &models.Pet{
					ID:      uuid.New(),
					StoreID: uuid.New(),
					Name:    "Fluffy",
					Status:  models.PetStatusAvailable,
				}
				repo.On("GetByID", mock.Anything, mock.Anything).Return(pet, nil)
				cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				repo.On("Delete", mock.Anything, mock.Anything).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockEncryptor := new(mocks.MockEncryptor)

			tt.setup(mockRepo, mockCache, mockEncryptor)

			service := NewPetService(mockRepo, mockCache, mockEncryptor)
			petID := uuid.New()

			err := service.DeletePetByID(context.Background(), petID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
			mockCache.AssertExpectations(t)
		})
	}
}

func TestPetService_MarkPetAsSold(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		setup   func(*mocks.MockPetRepository)
	}{
		{
			name:    "successful mark as sold",
			wantErr: false,
			setup: func(repo *mocks.MockPetRepository) {
				repo.On("Transaction", mock.AnythingOfType("func(*sql.Tx) error")).Return(nil)
			},
		},
		{
			name:    "transaction error",
			wantErr: true,
			setup: func(repo *mocks.MockPetRepository) {
				repo.On("Transaction", mock.AnythingOfType("func(*sql.Tx) error")).Return(assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockEncryptor := new(mocks.MockEncryptor)

			tt.setup(mockRepo)

			service := NewPetService(mockRepo, mockCache, mockEncryptor)
			petID := uuid.New()

			err := service.MarkPetAsSold(context.Background(), petID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPetService_DecryptBreederEmail(t *testing.T) {
	tests := []struct {
		name          string
		encryptedText string
		wantErr       bool
		setup         func(*mocks.MockEncryptor)
	}{
		{
			name:          "successful decryption",
			encryptedText: "encrypted_email",
			wantErr:       false,
			setup: func(encryptor *mocks.MockEncryptor) {
				encryptor.On("Decrypt", "encrypted_email").Return("john@example.com", nil)
			},
		},
		{
			name:          "empty encrypted email",
			encryptedText: "",
			wantErr:       true,
			setup:         func(*mocks.MockEncryptor) {}, // No mock needed
		},
		{
			name:          "decryption error",
			encryptedText: "invalid_encrypted",
			wantErr:       true,
			setup: func(encryptor *mocks.MockEncryptor) {
				encryptor.On("Decrypt", "invalid_encrypted").Return("", assert.AnError)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockPetRepository)
			mockCache := new(mocks.MockCache)
			mockEncryptor := new(mocks.MockEncryptor)

			tt.setup(mockEncryptor)

			service := NewPetService(mockRepo, mockCache, mockEncryptor)

			email, err := service.DecryptBreederEmail(tt.encryptedText)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, email)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "john@example.com", email)
			}

			mockEncryptor.AssertExpectations(t)
		})
	}
}

func TestPetServiceInterface_Implementation(t *testing.T) {
	// Test that PetService implements PetServiceInterface
	mockRepo := new(mocks.MockPetRepository)
	mockCache := new(mocks.MockCache)
	mockEncryptor := new(mocks.MockEncryptor)

	var _ PetServiceInterface = NewPetService(mockRepo, mockCache, mockEncryptor)
}
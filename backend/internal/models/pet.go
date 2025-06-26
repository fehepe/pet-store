package models

import (
	"time"

	"github.com/google/uuid"
)

type PetSpecies string

const (
	PetSpeciesCat  PetSpecies = "Cat"
	PetSpeciesDog  PetSpecies = "Dog"
	PetSpeciesFrog PetSpecies = "Frog"
)

type PetStatus string

const (
	PetStatusAvailable PetStatus = "available"
	PetStatusSold      PetStatus = "sold"
)

type Pet struct {
	ID                    uuid.UUID  `db:"id"`
	StoreID               uuid.UUID  `db:"store_id"`
	Name                  string     `db:"name"`
	Species               PetSpecies `db:"species"`
	Age                   int        `db:"age"`
	PictureURL            *string    `db:"picture_url"`
	Description           *string    `db:"description"`
	BreederName           string     `db:"breeder_name"`
	BreederEmailEncrypted string     `db:"breeder_email_encrypted"`
	Status                PetStatus  `db:"status"`
	CreatedAt             time.Time  `db:"created_at"`
	UpdatedAt             time.Time  `db:"updated_at"`
}

type CreatePetInput struct {
	StoreID      uuid.UUID
	Name         string
	Species      PetSpecies
	Age          int
	PictureURL   *string
	Description  *string
	BreederName  string
	BreederEmail string
}

type PetFilter struct {
	StoreID   *uuid.UUID
	Status    *PetStatus
	StartDate *time.Time
	EndDate   *time.Time
	Limit     int
	Offset    int
}

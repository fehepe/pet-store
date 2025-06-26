package models

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	OwnerID   string    `db:"owner_id"` // This will be the merchant's username
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type CreateStoreInput struct {
	Name    string
	OwnerID string
}


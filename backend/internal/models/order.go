package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID         uuid.UUID `db:"id"`
	CustomerID string    `db:"customer_id"`
	StoreID    uuid.UUID `db:"store_id"`
	TotalPets  int       `db:"total_pets"`
	CreatedAt  time.Time `db:"created_at"`
}

type OrderItem struct {
	ID          uuid.UUID `db:"id"`
	OrderID     uuid.UUID `db:"order_id"`
	PetID       uuid.UUID `db:"pet_id"`
	PurchasedAt time.Time `db:"purchased_at"`
}

type CreateOrderInput struct {
	CustomerID string
	StoreID    uuid.UUID
	PetIDs     []uuid.UUID
}

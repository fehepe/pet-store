package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fehepe/pet-store/backend/internal/database"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
)

// OrderRepositoryInterface defines the interface for order data operations
type OrderRepositoryInterface interface {
	CreateWithTx(ctx context.Context, tx *sql.Tx, order *models.Order) error
	CreateItem(ctx context.Context, tx *sql.Tx, item *models.OrderItem) error
	GetOrderPets(ctx context.Context, orderID uuid.UUID) ([]*models.Pet, error)
	UpdateWithTx(ctx context.Context, tx *sql.Tx, order *models.Order) error
	Transaction(fn func(*sql.Tx) error) error
}

// OrderRepository implements OrderRepositoryInterface
type OrderRepository struct {
	BaseRepository
}

// NewOrderRepository creates a new order repository
func NewOrderRepository(db database.Repository) OrderRepositoryInterface {
	return &OrderRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// CreateWithTx inserts a new order within a transaction
func (r *OrderRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, order *models.Order) error {
	query := `
		INSERT INTO orders (id, customer_id, store_id, total_pets, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`

	row := r.QueryInsertWithTx(ctx, tx, query,
		order.ID, order.CustomerID, order.StoreID, order.TotalPets, order.CreatedAt,
	)

	return row.Scan(
		&order.ID, &order.CustomerID, &order.StoreID, &order.TotalPets, &order.CreatedAt,
	)
}

// CreateItem inserts a new order item within a transaction
func (r *OrderRepository) CreateItem(ctx context.Context, tx *sql.Tx, item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (id, order_id, pet_id, purchased_at)
		VALUES ($1, $2, $3, $4)`

	return r.ExecInsertWithTx(ctx, tx, query,
		item.ID, item.OrderID, item.PetID, item.PurchasedAt,
	)
}

// GetOrderPets retrieves pets for a specific order
func (r *OrderRepository) GetOrderPets(ctx context.Context, orderID uuid.UUID) ([]*models.Pet, error) {
	query := `
		SELECT p.id, p.store_id, p.name, p.species, p.age, p.picture_url, p.description,
			   p.breeder_name, p.breeder_email_encrypted, p.status, p.created_at, p.updated_at
		FROM pets p
		JOIN order_items oi ON p.id = oi.pet_id
		WHERE oi.order_id = $1
		ORDER BY oi.purchased_at`

	rows, err := r.DB().QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order pets: %w", err)
	}
	defer rows.Close()

	pets := []*models.Pet{}
	for rows.Next() {
		var pet models.Pet
		err := rows.Scan(
			&pet.ID, &pet.StoreID, &pet.Name, &pet.Species, &pet.Age,
			&pet.PictureURL, &pet.Description, &pet.BreederName,
			&pet.BreederEmailEncrypted, &pet.Status, &pet.CreatedAt, &pet.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pet: %w", err)
		}
		pets = append(pets, &pet)
	}

	return pets, nil
}

// UpdateWithTx updates an existing order within a transaction
func (r *OrderRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, order *models.Order) error {
	query := `
		UPDATE orders 
		SET total_pets = $2
		WHERE id = $1
		RETURNING *`

	row := tx.QueryRowContext(ctx, query, order.ID, order.TotalPets)

	return row.Scan(
		&order.ID, &order.CustomerID, &order.StoreID, &order.TotalPets, &order.CreatedAt,
	)
}


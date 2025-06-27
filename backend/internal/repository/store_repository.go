package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fehepe/pet-store/backend/internal/database"
	"github.com/fehepe/pet-store/backend/internal/models"
)

// StoreRepositoryInterface defines the interface for store data operations
type StoreRepositoryInterface interface {
	Create(ctx context.Context, store *models.Store) error
	GetByOwnerID(ctx context.Context, ownerID string) (*models.Store, error)
	ListAll(ctx context.Context) ([]*models.Store, error)
}

// StoreRepository implements StoreRepositoryInterface
type StoreRepository struct {
	BaseRepository
}

// NewStoreRepository creates a new store repository
func NewStoreRepository(db database.Repository) StoreRepositoryInterface {
	return &StoreRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new store into the database
func (r *StoreRepository) Create(ctx context.Context, store *models.Store) error {
	query := `
		INSERT INTO stores (id, name, owner_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING *`

	row := r.QueryInsert(ctx, query,
		store.ID, store.Name, store.OwnerID, store.CreatedAt, store.UpdatedAt,
	)

	return row.Scan(
		&store.ID, &store.Name, &store.OwnerID, &store.CreatedAt, &store.UpdatedAt,
	)
}

// GetByOwnerID retrieves a store by its owner ID
func (r *StoreRepository) GetByOwnerID(ctx context.Context, ownerID string) (*models.Store, error) {
	query := `
		SELECT id, name, owner_id, created_at, updated_at
		FROM stores
		WHERE owner_id = $1`

	var store models.Store
	row := r.DB().QueryRowContext(ctx, query, ownerID)
	err := row.Scan(
		&store.ID, &store.Name, &store.OwnerID, &store.CreatedAt, &store.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, sql.ErrNoRows // Let the service handle auto-creation
	} else if err != nil {
		return nil, fmt.Errorf("failed to get store by owner: %w", err)
	}

	return &store, nil
}

// ListAll retrieves all stores
func (r *StoreRepository) ListAll(ctx context.Context) ([]*models.Store, error) {
	query := `
		SELECT id, name, owner_id, created_at, updated_at
		FROM stores
		ORDER BY name ASC`

	rows, err := r.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query stores: %w", err)
	}
	defer rows.Close()

	var stores []*models.Store
	for rows.Next() {
		var store models.Store
		err := rows.Scan(
			&store.ID, &store.Name, &store.OwnerID, &store.CreatedAt, &store.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan store: %w", err)
		}
		stores = append(stores, &store)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating store rows: %w", err)
	}

	return stores, nil
}

package repository

import (
	"context"
	"database/sql"

	"github.com/fehepe/pet-store/backend/internal/database"
)

// BaseRepository provides common repository functionality
type BaseRepository struct {
	db database.Repository
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db database.Repository) BaseRepository {
	return BaseRepository{db: db}
}

// Transaction executes a function within a database transaction
func (r *BaseRepository) Transaction(fn func(*sql.Tx) error) error {
	return r.db.Transaction(fn)
}

// DB returns the underlying database repository for direct access when needed
func (r *BaseRepository) DB() database.Repository {
	return r.db
}

// ExecInsertWithTx executes an INSERT query within a transaction using ExecContext
func (r *BaseRepository) ExecInsertWithTx(ctx context.Context, tx *sql.Tx, query string, args ...any) error {
	_, err := tx.ExecContext(ctx, query, args...)
	return err
}

// QueryInsertWithTx executes an INSERT query with RETURNING clause within a transaction
func (r *BaseRepository) QueryInsertWithTx(ctx context.Context, tx *sql.Tx, query string, args ...any) *sql.Row {
	return tx.QueryRowContext(ctx, query, args...)
}


// QueryInsert executes an INSERT query with RETURNING clause
func (r *BaseRepository) QueryInsert(ctx context.Context, query string, args ...any) *sql.Row {
	return r.db.QueryRowContext(ctx, query, args...)
}

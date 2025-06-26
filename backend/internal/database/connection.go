package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fehepe/pet-store/backend/internal/config"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Repository defines the interface for database operations
type Repository interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	Transaction(fn func(*sql.Tx) error) error

	Close() error
	Ping() error

	Migrate() error
}

// Ensure DB implements Repository interface
var _ Repository = (*DB)(nil)

type DB struct {
	*sql.DB
}

func NewConnection(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.DatabaseURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{db}, nil
}

func (db *DB) Migrate() error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	relativePaths := []string{
		"internal/database/migrations",
		"backend/internal/database/migrations",
		"../../internal/database/migrations",
	}

	var m *migrate.Migrate
	var lastErr error
	var foundPath string

	for _, relPath := range relativePaths {
		absPath := filepath.Join(cwd, relPath)

		if info, err := os.Stat(absPath); err == nil && info.IsDir() {
			foundPath = absPath
			fileURL := fmt.Sprintf("file://%s", absPath)

			m, err = migrate.NewWithDatabaseInstance(fileURL, "postgres", driver)
			if err == nil {
				break
			}
			lastErr = err
		}
	}

	if m == nil {
		return fmt.Errorf("failed to create migration instance from path %s: %w", foundPath, lastErr)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	fmt.Println("Database migrations completed successfully")
	return nil
}

func (db *DB) Transaction(fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

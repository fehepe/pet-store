package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fehepe/pet-store/backend/internal/database"
	apperrors "github.com/fehepe/pet-store/backend/internal/errors"
	"github.com/fehepe/pet-store/backend/internal/models"
	"github.com/google/uuid"
)

// PetRepositoryInterface defines the interface for pet data operations
type PetRepositoryInterface interface {
	Create(ctx context.Context, pet *models.Pet) error
	GetByID(ctx context.Context, petID uuid.UUID) (*models.Pet, error)
	List(ctx context.Context, filter models.PetFilter) ([]*models.Pet, int, error)
	Delete(ctx context.Context, petID uuid.UUID) error
	MarkAsSold(ctx context.Context, tx *sql.Tx, petID uuid.UUID) error
	Transaction(fn func(*sql.Tx) error) error
}

// PetRepository implements PetRepositoryInterface
type PetRepository struct {
	BaseRepository
}

// NewPetRepository creates a new pet repository
func NewPetRepository(db database.Repository) PetRepositoryInterface {
	return &PetRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create inserts a new pet into the database
func (r *PetRepository) Create(ctx context.Context, pet *models.Pet) error {
	query := `
		INSERT INTO pets (id, store_id, name, species, age, picture_url, description, 
			breeder_name, breeder_email_encrypted, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING *`

	row := r.QueryInsert(ctx, query,
		pet.ID, pet.StoreID, pet.Name, pet.Species, pet.Age,
		pet.PictureURL, pet.Description, pet.BreederName,
		pet.BreederEmailEncrypted, pet.Status, pet.CreatedAt, pet.UpdatedAt,
	)

	return row.Scan(
		&pet.ID, &pet.StoreID, &pet.Name, &pet.Species, &pet.Age,
		&pet.PictureURL, &pet.Description, &pet.BreederName,
		&pet.BreederEmailEncrypted, &pet.Status, &pet.CreatedAt, &pet.UpdatedAt,
	)
}

// GetByID retrieves a pet by its ID
func (r *PetRepository) GetByID(ctx context.Context, petID uuid.UUID) (*models.Pet, error) {
	query := `
		SELECT id, store_id, name, species, age, picture_url, description,
			   breeder_name, breeder_email_encrypted, status, created_at, updated_at
		FROM pets
		WHERE id = $1`

	var pet models.Pet
	row := r.DB().QueryRowContext(ctx, query, petID)
	err := row.Scan(
		&pet.ID, &pet.StoreID, &pet.Name, &pet.Species, &pet.Age,
		&pet.PictureURL, &pet.Description, &pet.BreederName,
		&pet.BreederEmailEncrypted, &pet.Status, &pet.CreatedAt, &pet.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, apperrors.NewPetNotFound(petID)
	} else if err != nil {
		return nil, fmt.Errorf("failed to get pet: %w", err)
	}

	return &pet, nil
}

// List retrieves pets with filtering and pagination
func (r *PetRepository) List(ctx context.Context, filter models.PetFilter) ([]*models.Pet, int, error) {
	var whereConditions []string
	var args []any
	argIndex := 1

	if filter.StoreID != nil && *filter.StoreID != uuid.Nil {
		whereConditions = append(whereConditions, fmt.Sprintf("store_id = $%d", argIndex))
		args = append(args, *filter.StoreID)
		argIndex++
	}
	if filter.Status != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, *filter.Status)
		argIndex++
	}
	if filter.StartDate != nil && filter.EndDate != nil {
		whereConditions = append(whereConditions, fmt.Sprintf("created_at BETWEEN $%d AND $%d", argIndex, argIndex+1))
		args = append(args, *filter.StartDate, *filter.EndDate)
		argIndex += 2
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM pets %s", whereClause)
	var total int
	err := r.DB().QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count pets: %w", err)
	}

	limit := 20
	offset := 0
	if filter.Limit > 0 {
		limit = filter.Limit
	}
	if filter.Offset >= 0 {
		offset = filter.Offset
	}

	// Add LIMIT and OFFSET parameters
	limitIndex := argIndex
	offsetIndex := argIndex + 1
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT id, store_id, name, species, age, picture_url, description,
			   breeder_name, breeder_email_encrypted, status, created_at, updated_at
		FROM pets
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, limitIndex, offsetIndex)

	rows, err := r.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query pets: %w", err)
	}
	defer rows.Close()

	var pets []*models.Pet
	for rows.Next() {
		var pet models.Pet
		err := rows.Scan(
			&pet.ID, &pet.StoreID, &pet.Name, &pet.Species, &pet.Age,
			&pet.PictureURL, &pet.Description, &pet.BreederName,
			&pet.BreederEmailEncrypted, &pet.Status, &pet.CreatedAt, &pet.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan pet: %w", err)
		}
		pets = append(pets, &pet)
	}

	return pets, total, nil
}

// Delete removes a pet from the database
func (r *PetRepository) Delete(ctx context.Context, petID uuid.UUID) error {
	query := `DELETE FROM pets WHERE id = $1`
	result, err := r.DB().ExecContext(ctx, query, petID)
	if err != nil {
		return fmt.Errorf("failed to delete pet: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperrors.NewPetNotFound(petID)
	}

	return nil
}

// MarkAsSold marks a pet as sold within a transaction
func (r *PetRepository) MarkAsSold(ctx context.Context, tx *sql.Tx, petID uuid.UUID) error {
	query := `UPDATE pets SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	result, err := tx.ExecContext(ctx, query, models.PetStatusSold, petID)
	if err != nil {
		return fmt.Errorf("failed to mark pet as sold: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return apperrors.NewPetNotFound(petID)
	}

	return nil
}

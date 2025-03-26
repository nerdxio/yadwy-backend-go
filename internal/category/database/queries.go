package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for categories
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new category repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create inserts a new category into the database
func (r *Repository) Create(ctx context.Context, name string, description string) (int64, error) {
	var id int64
	query := `
		INSERT INTO categories (name, description) 
		VALUES ($1, $2) 
		RETURNING id
	`

	descriptionSQL := sql.NullString{
		String: description,
		Valid:  description != "",
	}

	err := r.db.QueryRowxContext(ctx, query, name, descriptionSQL).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create category: %w", err)
	}

	return id, nil
}

// GetByID retrieves a category by its ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*Category, error) {
	var category Category
	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM categories 
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

// List retrieves all categories
func (r *Repository) List(ctx context.Context) ([]Category, error) {
	var categories []Category
	query := `
		SELECT id, name, description, created_at, updated_at 
		FROM categories 
		ORDER BY name
	`

	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

// Update updates an existing category
func (r *Repository) Update(ctx context.Context, id int64, name string, description string) error {
	query := `
		UPDATE categories 
		SET name = $2, description = $3, updated_at = NOW() 
		WHERE id = $1
	`

	descriptionSQL := sql.NullString{
		String: description,
		Valid:  description != "",
	}

	result, err := r.db.ExecContext(ctx, query, id, name, descriptionSQL)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// Delete removes a category from the database
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

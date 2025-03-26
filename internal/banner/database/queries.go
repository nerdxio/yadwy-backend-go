package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Repository handles database operations for banners
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new banner repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

// Create inserts a new banner into the database
func (r *Repository) Create(ctx context.Context, title, imageURL, targetURL, position string, isActive bool) (int64, error) {
	var id int64
	query := `
		INSERT INTO banners (title, image_url, target_url, is_active, position) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id
	`

	targetURLSQL := sql.NullString{
		String: targetURL,
		Valid:  targetURL != "",
	}

	err := r.db.QueryRowxContext(ctx, query, title, imageURL, targetURLSQL, isActive, position).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create banner: %w", err)
	}

	return id, nil
}

// GetByID retrieves a banner by its ID
func (r *Repository) GetByID(ctx context.Context, id int64) (*Banner, error) {
	var banner Banner
	query := `
		SELECT id, title, image_url, target_url, is_active, position, created_at, updated_at 
		FROM banners 
		WHERE id = $1
	`

	err := r.db.GetContext(ctx, &banner, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	return &banner, nil
}

// List retrieves all banners
func (r *Repository) List(ctx context.Context) ([]Banner, error) {
	var banners []Banner
	query := `
		SELECT id, title, image_url, target_url, is_active, position, created_at, updated_at 
		FROM banners 
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &banners, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}

	return banners, nil
}

// ListActive retrieves all active banners
func (r *Repository) ListActive(ctx context.Context) ([]Banner, error) {
	var banners []Banner
	query := `
		SELECT id, title, image_url, target_url, is_active, position, created_at, updated_at 
		FROM banners 
		WHERE is_active = true
		ORDER BY created_at DESC
	`

	err := r.db.SelectContext(ctx, &banners, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list active banners: %w", err)
	}

	return banners, nil
}

// Update updates an existing banner
func (r *Repository) Update(ctx context.Context, id int64, title, imageURL, targetURL, position string, isActive bool) error {
	query := `
		UPDATE banners 
		SET title = $2, image_url = $3, target_url = $4, is_active = $5, position = $6, updated_at = NOW() 
		WHERE id = $1
	`

	targetURLSQL := sql.NullString{
		String: targetURL,
		Valid:  targetURL != "",
	}

	result, err := r.db.ExecContext(ctx, query, id, title, imageURL, targetURLSQL, isActive, position)
	if err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("banner not found")
	}

	return nil
}

// Delete removes a banner from the database
func (r *Repository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM banners WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("banner not found")
	}

	return nil
}

package repository

import (
	"database/sql"

	"yadwy-backend/internal/category/domain"
)

// PostgresRepository implements the domain.CategoryRepository interface
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository creates a new postgres category repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// CreateCategory inserts a new category and returns its ID
// Implements the domain.CategoryRepository interface
func (r *PostgresRepository) CreateCategory(name, description string) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		name, description,
	).Scan(&id)

	return id, err
}

// Ensure PostgresRepository implements domain.CategoryRepository
var _ domain.CategoryRepository = (*PostgresRepository)(nil)

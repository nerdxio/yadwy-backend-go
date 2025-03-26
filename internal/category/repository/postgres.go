package repository

import (
	"context"
	"database/sql"

	"yadwy-backend/internal/category/domain"
	"yadwy-backend/sqlc/generated"
)

// PostgresRepository implements domain.CategoryRepository
type PostgresRepository struct {
	queries *generated.Queries
	db      *sql.DB
}

// NewPostgresRepository creates a new postgres repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		queries: generated.New(db),
		db:      db,
	}
}

// CreateCategory creates a new category
func (r *PostgresRepository) CreateCategory(name, description string) (int, error) {
	params := MapToDBParams(name, description)

	category, err := r.queries.CreateCategory(context.Background(), params)
	if err != nil {
		return 0, err
	}

	return int(category.ID), nil
}

// Ensure PostgresRepository implements domain.CategoryRepository
var _ domain.CategoryRepository = (*PostgresRepository)(nil)

package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"yadwy-backend/internal/category/database"
	"yadwy-backend/internal/category/domain"
)

// PostgresRepository implements domain.CategoryRepository
type PostgresRepository struct {
	repo *database.Repository
}

// NewPostgresRepository creates a new postgres repository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		repo: database.NewRepository(db),
	}
}

// CreateCategory creates a new category
func (r *PostgresRepository) CreateCategory(name, description string) (int, error) {
	id, err := r.repo.Create(context.Background(), name, description)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// GetCategory gets a category by ID
func (r *PostgresRepository) GetCategory(id int) (*domain.Category, error) {
	category, err := r.repo.GetByID(context.Background(), int64(id))
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, nil
	}

	model := category.ToModel()
	return &model, nil
}

// ListCategories lists all categories
func (r *PostgresRepository) ListCategories() ([]domain.Category, error) {
	categories, err := r.repo.List(context.Background())
	if err != nil {
		return nil, err
	}

	return database.ToModels(categories), nil
}

// UpdateCategory updates a category
func (r *PostgresRepository) UpdateCategory(id int, name, description string) error {
	return r.repo.Update(context.Background(), int64(id), name, description)
}

// DeleteCategory deletes a category
func (r *PostgresRepository) DeleteCategory(id int) error {
	return r.repo.Delete(context.Background(), int64(id))
}

// Ensure PostgresRepository implements domain.CategoryRepository
var _ domain.CategoryRepository = (*PostgresRepository)(nil)

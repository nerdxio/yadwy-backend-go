package service

import (
	"log/slog"

	"yadwy-backend/internal/category/domain"
)

// CategoryService handles category business logic
type CategoryService struct {
	repo domain.CategoryRepository
}

// NewCategoryService creates a new category service
func NewCategoryService(repo domain.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

// CreateCategory creates a new category
func (s *CategoryService) CreateCategory(name, description string) (int, error) {
	slog.Info("Creating category")

	id, err := s.repo.CreateCategory(name, description)
	if err != nil {
		slog.Error("Failed to create category", "error", err)
		return 0, err
	}

	return id, nil
}

// GetCategory retrieves a category by ID
func (s *CategoryService) GetCategory(id int) (*domain.Category, error) {
	slog.Info("Getting category", "id", id)

	category, err := s.repo.GetCategory(id)
	if err != nil {
		slog.Error("Failed to get category", "id", id, "error", err)
		return nil, err
	}

	return category, nil
}

// ListCategories retrieves all categories
func (s *CategoryService) ListCategories() ([]domain.Category, error) {
	slog.Info("Listing all categories")

	categories, err := s.repo.ListCategories()
	if err != nil {
		slog.Error("Failed to list categories", "error", err)
		return nil, err
	}

	return categories, nil
}

// UpdateCategory updates a category
func (s *CategoryService) UpdateCategory(id int, name, description string) error {
	slog.Info("Updating category", "id", id)

	err := s.repo.UpdateCategory(id, name, description)
	if err != nil {
		slog.Error("Failed to update category", "id", id, "error", err)
		return err
	}

	return nil
}

// DeleteCategory deletes a category
func (s *CategoryService) DeleteCategory(id int) error {
	slog.Info("Deleting category", "id", id)

	err := s.repo.DeleteCategory(id)
	if err != nil {
		slog.Error("Failed to delete category", "id", id, "error", err)
		return err
	}

	return nil
}

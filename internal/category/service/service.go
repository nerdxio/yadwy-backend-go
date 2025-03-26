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
func (s *CategoryService) CreateCategory(name, description string) (*domain.Category, error) {
	slog.Info("Creating category")

	category := domain.NewCategory(name, description)

	id, err := s.repo.CreateCategory(name, description)
	if err != nil {
		slog.Error("Failed to create category", "error", err)
		return nil, err
	}

	category.ID = id
	return &category, nil
}

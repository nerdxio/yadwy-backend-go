package category

import (
	"log/slog"
)

type CategoryService struct {
	repo CategoryRepo
}

func NewCategoryService(repo CategoryRepo) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}

func (uc *CategoryService) Execute(name, description string) (*Category, error) {
	slog.Info("Creating category")

	category := NewCategory(name, description)

	id, err := uc.repo.CreateCategory(name, description)
	if err != nil {
		slog.Error("Failed to create category", "error", err)
		return nil, err
	}

	category.ID = id
	return &category, nil
}

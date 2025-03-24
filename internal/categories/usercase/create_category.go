package usercase

import (
	"log/slog"
	"yadwy-backend/internal/categories/domain/contracts"
	"yadwy-backend/internal/categories/domain/modles"
)

type CreateCategoryUseCase struct {
	repo contracts.CategoryRepo
}

func NewCreateCategoryUseCase(repo contracts.CategoryRepo) *CreateCategoryUseCase {
	return &CreateCategoryUseCase{
		repo: repo,
	}
}

func (uc *CreateCategoryUseCase) Execute(name, description string) (*modles.Category, error) {
	slog.Info("Creating category")

	category := modles.NewCategory(name, description)

	id, err := uc.repo.CreateCategory(name, description)
	if err != nil {
		slog.Error("Failed to create category", "error", err)
		return nil, err
	}

	category.ID = id
	return &category, nil
}

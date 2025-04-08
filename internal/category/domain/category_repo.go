package domain

import (
	"context"
)

type CategoryRepo interface {
	CreateCategory(name, description, imageUrl string) (int, error)
	GetAllCategories(ctx context.Context) ([]Category, error)
}

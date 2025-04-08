package infra

import (
	"context"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
	"time"
	"yadwy-backend/internal/category/domain"
)

type CategoryDbo struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	ImageUrl    string    `db:"image_url"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
type CategoryRepoImpl struct {
	db     *sqlx.DB
	logger *zap.Logger
}

func NewCategoryRepo(db *sqlx.DB, logger *zap.Logger) *CategoryRepoImpl {
	return &CategoryRepoImpl{
		db:     db,
		logger: logger,
	}
}

func (r *CategoryRepoImpl) CreateCategory(name, description, imageUrl string) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO categories (name, description,image_url) VALUES ($1, $2, $3) RETURNING id",
		name, description, imageUrl,
	).Scan(&id)

	if err != nil {
		r.logger.Error("Failed to create category", zap.Error(err))
		return 0, err
	}
	return id, err
}

func (r *CategoryRepoImpl) GetAllCategories(ctx context.Context) ([]domain.Category, error) {
	var dbos []CategoryDbo
	err := r.db.SelectContext(ctx, &dbos, "SELECT * FROM categories")
	if err != nil {
		r.logger.Error("Failed to get categories", zap.Error(err))
		return nil, err
	}
	return mapToCategories(dbos), nil
}

func mapToCategories(dbos []CategoryDbo) []domain.Category {
	result := make([]domain.Category, 0, len(dbos))
	for _, dbo := range dbos {
		result = append(result, mapToCategory(dbo))
	}
	return result
}

func mapToCategory(dbo CategoryDbo) domain.Category {
	p := domain.Params{
		Id:          dbo.ID,
		Name:        dbo.Name,
		Description: dbo.Description,
		ImageUrl:    dbo.ImageUrl,
		CreatedAt:   dbo.CreatedAt,
		UpdatedAt:   dbo.UpdatedAt,
	}
	return domain.NewCategory(p)
}

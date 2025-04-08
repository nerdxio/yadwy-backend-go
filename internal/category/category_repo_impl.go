package category

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type CategoryDbo struct {
	ID          int    `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
	ImageUrl    string `db:"image_url"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
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
		"INSERT INTO categories (name, description,imageurl) VALUES ($1, $2, $3) RETURNING id",
		name, description, imageUrl,
	).Scan(&id)

	if err != nil {
		r.logger.Error("Failed to create category", zap.Error(err))
		return 0, err
	}
	return id, err
}

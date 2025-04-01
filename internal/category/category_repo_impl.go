package category

import (
	"github.com/jmoiron/sqlx"
)

type CategoryRepoImpl struct {
	db *sqlx.DB
}

func NewCategoryRepo(db *sqlx.DB) *CategoryRepoImpl {
	return &CategoryRepoImpl{
		db: db,
	}
}

func (r *CategoryRepoImpl) CreateCategory(name, description string) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO category (name, description) VALUES ($1, $2) RETURNING id",
		name, description,
	).Scan(&id)

	return id, err
}

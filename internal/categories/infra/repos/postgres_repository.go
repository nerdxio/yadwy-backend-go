package repos

import "database/sql"

type PostgresCategoryRepo struct {
	db *sql.DB
}

func NewPostgresCategoryRepo(db *sql.DB) *PostgresCategoryRepo {
	return &PostgresCategoryRepo{
		db: db,
	}
}

func (r *PostgresCategoryRepo) CreateCategory(name, description string) (int, error) {
	var id int
	err := r.db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id",
		name, description,
	).Scan(&id)

	return id, err
}

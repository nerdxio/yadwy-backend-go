package database

import (
	"database/sql"
	"time"
)

// Category represents a category in the database
type Category struct {
	ID          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

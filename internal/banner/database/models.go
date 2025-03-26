package database

import (
	"database/sql"
	"time"
)

// Banner represents a banner in the database
type Banner struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	ImageURL  string         `db:"image_url"`
	TargetURL sql.NullString `db:"target_url"`
	IsActive  bool           `db:"is_active"`
	Position  string         `db:"position"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

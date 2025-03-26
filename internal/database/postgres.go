package database

import (
	"database/sql"
	"fmt"
	"log/slog"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// Config holds database configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresDB creates a new PostgreSQL database connection
func NewPostgresDB(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}

	slog.Info("Connected to database", "host", cfg.Host, "database", cfg.DBName)
	return db, nil
}

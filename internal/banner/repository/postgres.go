package repository

import (
	"context"
	"database/sql"
	"fmt"

	"yadwy-backend/internal/banner/database/db"
	"yadwy-backend/internal/banner/domain"
)

// PostgresRepository implements domain.BannerRepository
type PostgresRepository struct {
	queries *db.Queries
	db      *sql.DB
}

// NewPostgresRepository creates a new postgres repository
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		queries: db.New(db),
		db:      db,
	}
}

// CreateBanner creates a new banner
func (r *PostgresRepository) CreateBanner(title, imageURL, targetURL, position string, isActive bool) (int, error) {
	params := MapToDBParams(title, imageURL, targetURL, position, isActive)

	banner, err := r.queries.CreateBanner(context.Background(), params)
	if err != nil {
		return 0, fmt.Errorf("failed to create banner: %w", err)
	}

	return int(banner.ID), nil
}

// GetBanner gets a banner by ID
func (r *PostgresRepository) GetBanner(id int) (*domain.Banner, error) {
	banner, err := r.queries.GetBanner(context.Background(), int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	result := MapToDomain(banner)
	return &result, nil
}

// ListBanners lists all banners
func (r *PostgresRepository) ListBanners() ([]domain.Banner, error) {
	banners, err := r.queries.ListBanners(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}

	result := make([]domain.Banner, len(banners))
	for i, banner := range banners {
		result[i] = MapToDomain(banner)
	}

	return result, nil
}

// ListActiveBanners lists all active banners
func (r *PostgresRepository) ListActiveBanners() ([]domain.Banner, error) {
	banners, err := r.queries.ListActiveBanners(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list active banners: %w", err)
	}

	result := make([]domain.Banner, len(banners))
	for i, banner := range banners {
		result[i] = MapToDomain(banner)
	}

	return result, nil
}

// UpdateBanner updates a banner
func (r *PostgresRepository) UpdateBanner(id int, title, imageURL, targetURL, position string, isActive bool) error {
	params := MapToUpdateParams(id, title, imageURL, targetURL, position, isActive)

	_, err := r.queries.UpdateBanner(context.Background(), params)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("banner not found")
		}
		return fmt.Errorf("failed to update banner: %w", err)
	}

	return nil
}

// DeleteBanner deletes a banner
func (r *PostgresRepository) DeleteBanner(id int) error {
	err := r.queries.DeleteBanner(context.Background(), int32(id))
	if err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}

	return nil
}

// Ensure PostgresRepository implements domain.BannerRepository
var _ domain.BannerRepository = (*PostgresRepository)(nil)

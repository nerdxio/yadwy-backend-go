package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"yadwy-backend/internal/banner/database"
	"yadwy-backend/internal/banner/domain"
)

// PostgresRepository implements domain.BannerRepository
type PostgresRepository struct {
	repo *database.Repository
}

// NewPostgresRepository creates a new postgres repository
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{
		repo: database.NewRepository(db),
	}
}

// CreateBanner creates a new banner
func (r *PostgresRepository) CreateBanner(title, imageURL, targetURL, position string, isActive bool) (int, error) {
	id, err := r.repo.Create(context.Background(), title, imageURL, targetURL, position, isActive)
	if err != nil {
		return 0, fmt.Errorf("failed to create banner: %w", err)
	}

	return int(id), nil
}

// GetBanner gets a banner by ID
func (r *PostgresRepository) GetBanner(id int) (*domain.Banner, error) {
	banner, err := r.repo.GetByID(context.Background(), int64(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get banner: %w", err)
	}

	if banner == nil {
		return nil, nil
	}

	result := banner.ToModel()
	return &result, nil
}

// ListBanners lists all banners
func (r *PostgresRepository) ListBanners() ([]domain.Banner, error) {
	banners, err := r.repo.List(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list banners: %w", err)
	}

	result := database.ToModels(banners)
	return result, nil
}

// ListActiveBanners lists all active banners
func (r *PostgresRepository) ListActiveBanners() ([]domain.Banner, error) {
	banners, err := r.repo.ListActive(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to list active banners: %w", err)
	}

	result := database.ToModels(banners)
	return result, nil
}

// UpdateBanner updates a banner
func (r *PostgresRepository) UpdateBanner(id int, title, imageURL, targetURL, position string, isActive bool) error {
	err := r.repo.Update(context.Background(), int64(id), title, imageURL, targetURL, position, isActive)
	if err != nil {
		return fmt.Errorf("failed to update banner: %w", err)
	}

	return nil
}

// DeleteBanner deletes a banner
func (r *PostgresRepository) DeleteBanner(id int) error {
	err := r.repo.Delete(context.Background(), int64(id))
	if err != nil {
		return fmt.Errorf("failed to delete banner: %w", err)
	}

	return nil
}

// Ensure PostgresRepository implements domain.BannerRepository
var _ domain.BannerRepository = (*PostgresRepository)(nil)

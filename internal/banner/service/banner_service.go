package service

import (
	"yadwy-backend/internal/banner/domain"
)

// BannerService handles banner business logic
type BannerService struct {
	repo domain.BannerRepository
}

// NewBannerService creates a new banner service
func NewBannerService(repo domain.BannerRepository) *BannerService {
	return &BannerService{
		repo: repo,
	}
}

// CreateBanner creates a new banner
func (s *BannerService) CreateBanner(title, imageURL, targetURL, position string, isActive bool) (int, error) {
	return s.repo.CreateBanner(title, imageURL, targetURL, position, isActive)
}

// GetBanner gets a banner by ID
func (s *BannerService) GetBanner(id int) (*domain.Banner, error) {
	return s.repo.GetBanner(id)
}

// ListBanners lists all banners
func (s *BannerService) ListBanners() ([]domain.Banner, error) {
	return s.repo.ListBanners()
}

// ListActiveBanners lists all active banners
func (s *BannerService) ListActiveBanners() ([]domain.Banner, error) {
	return s.repo.ListActiveBanners()
}

// UpdateBanner updates a banner
func (s *BannerService) UpdateBanner(id int, title, imageURL, targetURL, position string, isActive bool) error {
	return s.repo.UpdateBanner(id, title, imageURL, targetURL, position, isActive)
}

// DeleteBanner deletes a banner
func (s *BannerService) DeleteBanner(id int) error {
	return s.repo.DeleteBanner(id)
}

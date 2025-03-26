package domain

// BannerRepository defines the interface for banner data operations
type BannerRepository interface {
	CreateBanner(title, imageURL, targetURL, position string, isActive bool) (int, error)
	GetBanner(id int) (*Banner, error)
	ListBanners() ([]Banner, error)
	ListActiveBanners() ([]Banner, error)
	UpdateBanner(id int, title, imageURL, targetURL, position string, isActive bool) error
	DeleteBanner(id int) error
}

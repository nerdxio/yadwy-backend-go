package domain

import "time"

// Banner represents a banner in the system
type Banner struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	TargetURL string    `json:"target_url,omitempty"`
	IsActive  bool      `json:"is_active"`
	Position  string    `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBanner creates a new banner
func NewBanner(title, imageURL, targetURL, position string, isActive bool) Banner {
	return Banner{
		Title:     title,
		ImageURL:  imageURL,
		TargetURL: targetURL,
		IsActive:  isActive,
		Position:  position,
	}
}

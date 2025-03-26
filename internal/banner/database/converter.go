package database

import (
	"yadwy-backend/internal/banner/domain"
)

// ToModel converts a database banner to a domain model
func (b *Banner) ToModel() domain.Banner {
	targetURL := ""
	if b.TargetURL.Valid {
		targetURL = b.TargetURL.String
	}

	return domain.Banner{
		ID:        int(b.ID),
		Title:     b.Title,
		ImageURL:  b.ImageURL,
		TargetURL: targetURL,
		IsActive:  b.IsActive,
		Position:  b.Position,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

// ToModels converts a slice of database banners to domain models
func ToModels(banners []Banner) []domain.Banner {
	result := make([]domain.Banner, len(banners))
	for i, banner := range banners {
		result[i] = banner.ToModel()
	}
	return result
}

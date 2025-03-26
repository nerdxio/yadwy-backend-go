package repository

import (
	"database/sql"
	"yadwy-backend/internal/banner/database/db"
	"yadwy-backend/internal/banner/domain"
)

// MapToDomain converts a database model to a domain model
func MapToDomain(dbBanner db.Banner) domain.Banner {
	targetURL := ""
	if dbBanner.TargetUrl.Valid {
		targetURL = dbBanner.TargetUrl.String
	}

	return domain.Banner{
		ID:        int(dbBanner.ID),
		Title:     dbBanner.Title,
		ImageURL:  dbBanner.ImageUrl,
		TargetURL: targetURL,
		IsActive:  dbBanner.IsActive,
		Position:  dbBanner.Position,
		CreatedAt: dbBanner.CreatedAt,
		UpdatedAt: dbBanner.UpdatedAt,
	}
}

// MapToDBParams converts domain data to database parameters
func MapToDBParams(title, imageURL, targetURL, position string, isActive bool) db.CreateBannerParams {
	return db.CreateBannerParams{
		Title:    title,
		ImageUrl: imageURL,
		TargetUrl: sql.NullString{
			String: targetURL,
			Valid:  targetURL != "",
		},
		IsActive: isActive,
		Position: position,
	}
}

// MapToUpdateParams converts domain data to database update parameters
func MapToUpdateParams(id int, title, imageURL, targetURL, position string, isActive bool) db.UpdateBannerParams {
	return db.UpdateBannerParams{
		ID:       int32(id),
		Title:    title,
		ImageUrl: imageURL,
		TargetUrl: sql.NullString{
			String: targetURL,
			Valid:  targetURL != "",
		},
		IsActive: isActive,
		Position: position,
	}
}

package repository

import (
	"database/sql"
	"yadwy-backend/internal/category/domain"
	"yadwy-backend/sqlc/generated"
)

// MapToDomain converts a database model to a domain model
func MapToDomain(dbCategory generated.Category) domain.Category {
	description := ""
	if dbCategory.Description.Valid {
		description = dbCategory.Description.String
	}

	return domain.Category{
		ID:          int(dbCategory.ID),
		Name:        dbCategory.Name,
		Description: description,
		CreatedAt:   dbCategory.CreatedAt,
		UpdatedAt:   dbCategory.UpdatedAt,
	}
}

// MapToDBParams converts domain data to database parameters
func MapToDBParams(name, description string) generated.CreateCategoryParams {
	return generated.CreateCategoryParams{
		Name: name,
		Description: sql.NullString{
			String: description,
			Valid:  description != "",
		},
	}
}

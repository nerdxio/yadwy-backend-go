package database

import (
	"yadwy-backend/internal/category/domain"
)

// ToModel converts a database category to a domain model
func (c *Category) ToModel() domain.Category {
	description := ""
	if c.Description.Valid {
		description = c.Description.String
	}

	return domain.Category{
		ID:          int(c.ID),
		Name:        c.Name,
		Description: description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// ToModels converts a slice of database categories to domain models
func ToModels(categories []Category) []domain.Category {
	result := make([]domain.Category, len(categories))
	for i, category := range categories {
		result[i] = category.ToModel()
	}
	return result
}

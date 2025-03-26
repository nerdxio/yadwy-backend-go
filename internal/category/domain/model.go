package domain

import "time"

// Category represents a product category
type Category struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewCategory creates a new category
func NewCategory(name, description string) Category {
	return Category{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

package modles

import "time"

// Category Domain model
type Category struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCategory(
	name string,
	description string,
) Category {
	return Category{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

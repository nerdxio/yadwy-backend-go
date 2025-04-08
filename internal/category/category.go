package category

import "time"

// Category Domain model
type Category struct {
	id          int64
	name        string
	description string
	imageUrl    string
	createdAt   time.Time
	updatedAt   time.Time
}

func NewCategory(
	name string,
	description string,
	imageUrl string,
) Category {
	return Category{
		name:        name,
		description: description,
		imageUrl:    imageUrl,
	}
}

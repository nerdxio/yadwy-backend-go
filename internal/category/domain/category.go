package domain

import "time"

type Category struct {
	id          int64
	name        string
	description string
	imageUrl    string
	createdAt   time.Time
	updatedAt   time.Time
}

type Params struct {
	Id          int64
	Name        string
	Description string
	ImageUrl    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewCategory(
	p Params,
) Category {
	return Category{
		id:          p.Id,
		name:        p.Name,
		description: p.Description,
		imageUrl:    p.ImageUrl,
		createdAt:   p.CreatedAt,
		updatedAt:   p.UpdatedAt,
	}
}

func (c *Category) Id() int64 {
	return c.id
}

func (c *Category) Name() string {
	return c.name
}

func (c *Category) Description() string {
	return c.description
}

func (c *Category) ImageUrl() string {
	return c.imageUrl
}

func (c *Category) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Category) UpdatedAt() time.Time {
	return c.updatedAt
}

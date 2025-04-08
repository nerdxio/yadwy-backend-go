package application

import (
	"time"
	"yadwy-backend/internal/category/domain"
)

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CategoryRes struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageUrl    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func MapToCategoryRes(c domain.Category) CategoryRes {
	return CategoryRes{
		Id:          c.Id(),
		Name:        c.Name(),
		Description: c.Description(),
		ImageUrl:    c.ImageUrl(),
		CreatedAt:   c.CreatedAt(),
		UpdatedAt:   c.UpdatedAt(),
	}
}

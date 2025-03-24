package handler

import (
	"github.com/go-chi/chi/v5"
)

type CategoryRoutes struct {
	categoryHandler *CreateCategoryHandler
}

func NewCategoryRoutes(categoryHandler *CreateCategoryHandler) *CategoryRoutes {
	return &CategoryRoutes{
		categoryHandler: categoryHandler,
	}
}

func (cr *CategoryRoutes) RegisterRoutes(r chi.Router) {
	r.Route("/categories", func(r chi.Router) {
		r.Post("/", cr.categoryHandler.Create)
		// Add more category endpoints here:
		// r.Get("/", cr.categoryHandler.List)
		// r.Get("/{id}", cr.categoryHandler.GetByID)
	})
}

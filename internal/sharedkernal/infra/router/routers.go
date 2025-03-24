package router

import (
	"database/sql"
	"fmt"
	"net/http"
	"yadwy-backend/internal/categories/infra/handler"
	"yadwy-backend/internal/categories/infra/repos"
	"yadwy-backend/internal/categories/usercase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func LoadRouters(db *sql.DB) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
		w.WriteHeader(http.StatusOK)
	})

	// sub routes
	router.Route("/categories", func(r chi.Router) {
		loadCategoryRoutes(db, r)
	})

	return router
}

func loadCategoryRoutes(db *sql.DB, r chi.Router) {
	categoryRepo := repos.NewPostgresCategoryRepo(db)
	createCategoryUseCase := usercase.NewCreateCategoryUseCase(categoryRepo)
	categoryHandler := handler.NewCategoryHandler(createCategoryUseCase)
	r.Route("/categories", func(r chi.Router) {
		r.Post("/", categoryHandler.Create)
		// Add more category endpoints here:
	})
}

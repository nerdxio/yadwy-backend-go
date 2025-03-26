package app

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"yadwy-backend/internal/category/api"
	"yadwy-backend/internal/category/repository"
	"yadwy-backend/internal/category/service"
)

// SetupRouter configures the Chi router with all application routes
func SetupRouter(db *sql.DB) http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	// Home route
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
		w.WriteHeader(http.StatusOK)
	})

	// API routes
	router.Route("/categories", func(r chi.Router) {
		loadCategoryRoutes(db, r)
	})

	return router
}

// loadCategoryRoutes sets up all category-related routes
func loadCategoryRoutes(db *sql.DB, r chi.Router) {
	// Initialize repository
	categoryRepo := repository.NewPostgresRepository(db)

	// Initialize service
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize handler
	categoryHandler := api.NewCategoryHandler(categoryService)

	// Register routes
	r.Post("/", categoryHandler.Create)
	// Add more category endpoints here
}

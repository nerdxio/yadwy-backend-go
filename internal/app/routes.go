package app

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"

	bannerAPI "yadwy-backend/internal/banner/api"
	"yadwy-backend/internal/banner/repository"
	bannerService "yadwy-backend/internal/banner/service"
	"yadwy-backend/internal/category/api"
	categoryRepo "yadwy-backend/internal/category/repository"
	"yadwy-backend/internal/category/service"
)

// SetupRouter configures the Chi router with all application routes
func SetupRouter(db *sqlx.DB) http.Handler {
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

	router.Route("/banners", func(r chi.Router) {
		loadBannerRoutes(db, r)
	})

	return router
}

// loadCategoryRoutes sets up all category-related routes
func loadCategoryRoutes(db *sqlx.DB, r chi.Router) {
	// Initialize repository
	categoryRepo := categoryRepo.NewPostgresRepository(db)

	// Initialize service
	categoryService := service.NewCategoryService(categoryRepo)

	// Initialize handler
	categoryHandler := api.NewCategoryHandler(categoryService)

	// Register routes
	r.Get("/", categoryHandler.ListCategories)
	r.Post("/", categoryHandler.CreateCategory)
}

// loadBannerRoutes sets up all banner-related routes
func loadBannerRoutes(db *sqlx.DB, r chi.Router) {
	// Initialize repository
	bannerRepo := repository.NewPostgresRepository(db)

	// Initialize service
	bannerSvc := bannerService.NewBannerService(bannerRepo)

	// Initialize handler
	bannerHandler := bannerAPI.NewBannerHandler(bannerSvc)

	// Register routes
	r.Get("/", bannerHandler.ListBanners)
	r.Get("/active", bannerHandler.ListActiveBanners)
	r.Post("/", bannerHandler.CreateBanner)
	r.Get("/{id}", bannerHandler.GetBanner)
	r.Put("/{id}", bannerHandler.UpdateBanner)
	r.Delete("/{id}", bannerHandler.DeleteBanner)
}

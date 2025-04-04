package app

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"net/http"
	ch "yadwy-backend/internal/category"
	uh "yadwy-backend/internal/users/handlers"
)

func SetupRouter(db *sqlx.DB, key string) http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Heartbeat("/ping"))

	router.Route("/users", func(r chi.Router) {
		uh.LoadUserRoutes(db, r, key)
	})

	router.Mount("/category", ch.LoadCategoryRoutes(db))
	return router
}

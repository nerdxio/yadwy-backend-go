package app

import (
	"net/http"
	"yadwy-backend/internal/users/application"
	"yadwy-backend/internal/users/db"
	"yadwy-backend/internal/users/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(db *sqlx.DB) http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Heartbeat("/ping"))

	router.Route("/users", func(r chi.Router) {
		loadUserRoutes(db, r)
	})

	return router
}

func loadUserRoutes(b *sqlx.DB, r chi.Router) {
	userRepo := db.NewUserRepo(b)
	userSvc := application.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userSvc)

	r.Post("/", userHandler.RegisterUser)
}

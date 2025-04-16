package app

import (
	"net/http"
	bh "yadwy-backend/internal/banner"
	carth "yadwy-backend/internal/cart/infra"
	ch "yadwy-backend/internal/category/infra"
	"yadwy-backend/internal/common"
	ph "yadwy-backend/internal/prodcuts/infra"
	uh "yadwy-backend/internal/users/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func SetupRouter(db *sqlx.DB, jwt *common.JWTGenerator, logger *zap.Logger) http.Handler {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Heartbeat("/ping"))

	// Routes to handle static files
	router.HandleFunc("/images/{image}", func(w http.ResponseWriter, r *http.Request) {
		chi.URLParam(r, "image")
		http.ServeFile(w, r, "/home/nerd/images/"+chi.URLParam(r, "image"))
	})

	router.Route("/users", func(r chi.Router) {
		uh.LoadUserRoutes(db, r, jwt)
	})

	router.Mount("/category", ch.LoadCategoryRoutes(db, logger, jwt))
	router.Mount("/banners", bh.LoadBannerRoutes(db, logger, jwt))
	router.Mount("/products", ph.LoadProductsRoutes(db, logger, jwt))
	router.Mount("/cart", carth.LoadCartRoutes(db, logger, jwt))
	return router
}

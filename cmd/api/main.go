// Package main Yadwy Backend API
// @title Yadwy Backend API
// @version 1.0
// @description This is the Yadwy backend service API documentation
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@yadwy.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer: ` prefix, e.g. "Bearer abcde12345".

package main

import (
	"context"
	"os"
	"os/signal"
	_ "yadwy-backend/api/swagger"
	"yadwy-backend/internal/common"

	"go.uber.org/zap"

	"yadwy-backend/internal/app"
	"yadwy-backend/internal/config"
	"yadwy-backend/internal/database"
)

func main() {
	logger, err := common.NewLogger()
	if err != nil {
		panic(err)
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load config", zap.Error(err))
		os.Exit(1)
	}

	if err := database.RunMigrations(cfg.Database); err != nil {
		logger.Error("Failed to run migrations", zap.Error(err))
		os.Exit(1)
	}

	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		os.Exit(1)
	}

	application := app.New(cfg, db, logger)

	application.Router = app.SetupRouter(db, application.JWT, application.Logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := application.Start(ctx); err != nil {
		logger.Error("Application error", zap.Error(err))
		os.Exit(1)
	}
}

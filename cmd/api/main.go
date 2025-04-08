package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"yadwy-backend/internal/common"

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

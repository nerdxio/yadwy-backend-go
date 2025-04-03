package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"yadwy-backend/internal/app"
	"yadwy-backend/internal/config"
	"yadwy-backend/internal/database"
)

func main() {
	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	if err := database.RunMigrations(cfg.Database); err != nil {
		slog.Error("Failed to run database migrations", "error", err)
		os.Exit(1)
	}

	// Connect to database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}

	application := app.New(cfg, db)

	// Setup router
	application.Router = app.SetupRouter(db, cfg.JWT.Secret)

	// Setup graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Start the application
	if err := application.Start(ctx); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}

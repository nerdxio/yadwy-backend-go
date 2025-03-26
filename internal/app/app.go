package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"yadwy-backend/internal/config"
)

// App represents the application
type App struct {
	Router http.Handler
	DB     *sqlx.DB
	config *config.Config
}

// New creates a new application instance
func New(cfg *config.Config, db *sqlx.DB) *App {
	return &App{
		DB:     db,
		config: cfg,
	}
}

// Start starts the HTTP server and handles graceful shutdown
func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.Server.Port),
		Handler: a.Router,
	}

	defer func() {
		if err := a.DB.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}()

	slog.Info("Starting server", "port", a.config.Server.Port)

	ch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		slog.Info("Shutting down server...")
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}

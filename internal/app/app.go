package app

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log/slog"
	"net/http"
	"time"
	"yadwy-backend/internal/common"

	"github.com/jmoiron/sqlx"
	"yadwy-backend/internal/config"
)

type App struct {
	Router http.Handler
	DB     *sqlx.DB
	Config *config.Config
	Logger *zap.Logger
	JWT    *common.JWTGenerator
}

func New(cfg *config.Config, db *sqlx.DB, logger *zap.Logger) *App {
	jwt := common.NewJWTGenerator(cfg.JWT.Secret)
	return &App{
		DB:     db,
		Config: cfg,
		Logger: logger,
		JWT:    jwt,
	}
}

// Start starts the HTTP server and handles graceful shutdown
func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.Config.Server.Port),
		Handler: a.Router,
	}

	defer func() {
		if err := a.DB.Close(); err != nil {
			slog.Error("Failed to close database connection", "error", err)
		}
	}()

	slog.Info("Starting server", "port", a.Config.Server.Port)

	ch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

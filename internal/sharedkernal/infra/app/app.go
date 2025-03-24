package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"yadwy-backend/config"
	"yadwy-backend/internal/sharedkernal/infra/database"
)

type App struct {
	Router http.Handler
	DB     *sql.DB
	config config.Config
}

func New(cfg config.Config) (*App, error) {
	db, err := database.NewPostgresConnection(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return &App{
		DB:     db,
		config: cfg,
	}, nil
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.config.ServerPort),
		Handler: a.Router,
	}

	defer func() {
		if err := a.DB.Close(); err != nil {
			log.Println("Failed to close database connection:", err)
		}
	}()

	log.Println("Starting the Server on port", server.Addr)

	ch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}
}

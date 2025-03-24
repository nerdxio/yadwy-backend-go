package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"yadwy-backend/config"
	"yadwy-backend/internal/sharedkernal/infra/app"
	"yadwy-backend/internal/sharedkernal/infra/router"
)

func main() {
	cfg := config.LoadConfig()

	app, err := app.New(cfg)
	if err != nil {
		log.Fatal("Failed to initialize application:", err)
	}

	app.Router = router.LoadRouters(app.DB)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err = app.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

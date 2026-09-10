package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/techcollectiveassam/flood-response-api/internal/api"
	"github.com/techcollectiveassam/flood-response-api/internal/app"
)

func main() {
	_ = godotenv.Load()
	application, err := app.New()
	if err != nil {
		log.Fatal(err)
	}
	defer application.Close()

	features := app.NewFeatures(application)

	router := api.NewRouter(features)

	server := &http.Server{
		Addr:    ":" + application.Config.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Printf("server listening on :%s", application.Config.Port)

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}

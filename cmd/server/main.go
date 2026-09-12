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

	logger := application.Logger

	features := app.NewFeatures(application)

	router := api.NewRouter(features, logger)

	server := &http.Server{
		Addr:    ":" + application.Config.Port,
		Handler: router,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("server listening", "port", application.Config.Port)

	<-ctx.Done()

	logger.Info("shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}

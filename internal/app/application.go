package app

import (
	"database/sql"
	"log"
	"log/slog"

	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
)

type Application struct {
	Config *config.Config
	DB     *sql.DB
	Logger *slog.Logger
}

func New() (*Application, error) {
	configuration, err := config.Load()
	if err != nil {
		return nil, err
	}

	db, err := database.New(configuration.Database)
	if err != nil {
		return nil, err
	}

	logger := slog.New(slog.NewJSONHandler(log.Writer(), &slog.HandlerOptions{
		AddSource: true,
	}))

	return &Application{
		Config: configuration,
		DB:     db,
		Logger: logger,
	}, nil
}

func (a *Application) Close() error {
	if a.DB == nil {
		return nil
	}

	return a.DB.Close()
}

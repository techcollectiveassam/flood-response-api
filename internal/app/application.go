package app

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/logging"
)

type Application struct {
	Config *config.Config
	DB     *pgxpool.Pool
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

	return &Application{
		Config: configuration,
		DB:     db,
		Logger: logging.New(configuration.Logger),
	}, nil
}

func (a *Application) Close() error {
	if a.DB == nil {
		return nil
	}

	a.DB.Close()
	return nil
}

package app

import (
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
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

	a.DB.Close()
	return nil
}

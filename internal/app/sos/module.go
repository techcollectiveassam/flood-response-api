package sos

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

type Module struct {
	handler *Handler
}

func New(db *pgxpool.Pool, cfg *config.Config) *Module {
	repository := NewRepository(db)
	service := NewService(repository, cfg.Sos.DuplicateRadiusMeters)
	handler := NewHandler(service)
	return &Module{
		handler: handler,
	}
}

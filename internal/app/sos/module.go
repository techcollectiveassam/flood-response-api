package sos

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
)

type Module struct {
	handler *Handler
}

func New(db *pgxpool.Pool, _ *config.Config) *Module {
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)
	return &Module{
		handler: handler,
	}
}

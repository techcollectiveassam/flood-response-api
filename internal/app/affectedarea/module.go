package affectedarea

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *Handler
}

func New(db *pgxpool.Pool, logger *slog.Logger) *Module {
	repository := NewRepository(db)
	service := NewService(repository, NewAffectedAreaResolver(), logger)
	handler := NewHandler(service)
	return &Module{handler: handler}
}

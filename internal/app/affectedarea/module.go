package affectedarea

import (
	"database/sql"
	"log/slog"
)

type Module struct {
	handler *Handler
}

func New(db *sql.DB, logger *slog.Logger) *Module {
	repository := NewRepository(db)
	service := NewService(repository, logger)
	handler := NewHandler(service)
	return &Module{handler: handler}
}

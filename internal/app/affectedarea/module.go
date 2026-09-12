package affectedarea

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Module struct {
	handler *Handler
}

func New(db *pgxpool.Pool) *Module {
	repository := NewRepository(db)
	service := NewService(repository, NewAffectedAreaResolver())
	handler := NewHandler(service)
	return &Module{handler: handler}
}

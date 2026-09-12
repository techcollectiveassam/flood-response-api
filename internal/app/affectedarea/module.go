package affectedarea

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/config"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/pagination"
)

type Module struct {
	handler *Handler
	params  pagination.Params
}

func New(db *pgxpool.Pool, cfg *config.Config) *Module {
	repository := NewRepository(db)
	service := NewService(repository, NewAffectedAreaResolver())
	handler := NewHandler(service)
	return &Module{
		handler: handler,
		params:  pagination.ParamsFromConfig(cfg),
	}
}

package affectedarea

import (
	"context"
	"database/sql"
	"errors"
)

var ErrNotImplemented = errors.New("affected area persistence is not implemented")

type Repository interface {
	Create(ctx context.Context, area *AffectedArea) error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	area *AffectedArea,
) error {
	// Query implementation intentionally belongs here, but is outside this skeleton.
	return ErrNotImplemented
}

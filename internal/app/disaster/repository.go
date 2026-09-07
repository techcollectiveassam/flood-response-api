package disaster

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

type Repository interface {
	Create(ctx context.Context, d *Disaster) error
	List(ctx context.Context) ([]Disaster, error)
}

type PostgresRepository struct {
	queries *sqlcgen.Queries
}

func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlcgen.New(pool),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, d *Disaster) error {
	result, err := r.queries.CreateDisaster(ctx, sqlcgen.CreateDisasterParams{
		Name:        d.Name,
		Description: database.NullableString(d.Description),
		Type:        sqlcgen.DisasterType(d.Type),
		Status:      sqlcgen.DisasterStatus(d.Status),
		Column5:     database.Timestamptz(d.StartsAt),
		EndsAt:      database.Timestamptz(d.EndsAt),
	})
	if err != nil {
		return database.TranslateError(err)
	}

	d.ID = result.ID
	d.Name = result.Name
	d.Description = database.TextValue(result.Description)
	d.Type = string(result.Type)
	d.Status = string(result.Status)
	d.StartsAt = database.FromTimestamptz(result.StartsAt)
	d.EndsAt = database.FromTimestamptz(result.EndsAt)
	return nil
}

func (r *PostgresRepository) List(ctx context.Context) ([]Disaster, error) {
	rows, err := r.queries.ListDisasters(ctx)
	if err != nil {
		return nil, database.TranslateError(err)
	}

	disasters := make([]Disaster, 0, len(rows))
	for _, row := range rows {
		disasters = append(disasters, Disaster{
			ID:          row.ID,
			Name:        row.Name,
			Description: database.TextValue(row.Description),
			Type:        string(row.Type),
			Status:      string(row.Status),
			StartsAt:    database.FromTimestamptz(row.StartsAt),
			EndsAt:      database.FromTimestamptz(row.EndsAt),
		})
	}
	return disasters, nil
}

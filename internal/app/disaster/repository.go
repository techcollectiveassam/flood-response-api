package disaster

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

var ErrDisasterNotFound = apperror.NotFound("disaster_not_found", "disaster not found")

type Repository interface {
	Create(ctx context.Context, d *Disaster) error
	List(ctx context.Context) ([]Disaster, error)
	GetByRefID(ctx context.Context, refID uuid.UUID) (*Disaster, error)
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
	d.RefID = result.RefID
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
		disasters = append(disasters, disasterFromRow(row.ID, row.RefID, row.Name, row.Description, row.Type, row.Status, row.StartsAt, row.EndsAt))
	}
	return disasters, nil
}

func (r *PostgresRepository) GetByRefID(ctx context.Context, refID uuid.UUID) (*Disaster, error) {
	row, err := r.queries.GetDisasterByRefID(ctx, refID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrDisasterNotFound
		}
		return nil, database.TranslateError(err)
	}

	d := disasterFromRow(row.ID, row.RefID, row.Name, row.Description, row.Type, row.Status, row.StartsAt, row.EndsAt)
	return &d, nil
}

func disasterFromRow(id int32, refID uuid.UUID, name string, description *string, dtype sqlcgen.DisasterType, status sqlcgen.DisasterStatus, startsAt, endsAt pgtype.Timestamptz) Disaster {
	return Disaster{
		ID:          id,
		RefID:       refID,
		Name:        name,
		Description: database.TextValue(description),
		Type:        string(dtype),
		Status:      string(status),
		StartsAt:    database.FromTimestamptz(startsAt),
		EndsAt:      database.FromTimestamptz(endsAt),
	}
}

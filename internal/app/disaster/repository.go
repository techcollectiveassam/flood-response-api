package disaster

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

type Repository interface {
	Create(ctx context.Context, d *Disaster) error
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
		Description: nullableString(d.Description),
		Type:        sqlcgen.DisasterType(d.Type),
		Status:      sqlcgen.DisasterStatus(d.Status),
		Column5:     timestamp(d.StartsAt),
		EndsAt:      toTimestamptz(d.EndsAt),
	})
	if err != nil {
		return database.TranslateError(err)
	}

	d.ID = result.ID
	d.Name = result.Name
	d.Description = textValue(result.Description)
	d.Type = string(result.Type)
	d.Status = string(result.Status)
	d.StartsAt = fromTimestamptz(result.StartsAt)
	d.EndsAt = fromTimestamptz(result.EndsAt)
	return nil
}

func nullableString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func textValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func timestamp(value *time.Time) interface{} {
	if value == nil || value.IsZero() {
		return nil
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func toTimestamptz(value *time.Time) pgtype.Timestamptz {
	if value == nil || value.IsZero() {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *value, Valid: true}
}

func fromTimestamptz(value pgtype.Timestamptz) *time.Time {
	if !value.Valid {
		return nil
	}
	t := value.Time
	return &t
}

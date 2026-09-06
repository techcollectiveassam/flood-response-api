package affectedarea

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

var ErrDisasterNotFound = apperror.NotFound("disaster_not_found", "disaster not found")

type Repository interface {
	Create(ctx context.Context, area *AffectedArea) error
}

type PostgresRepository struct {
	queries *sqlcgen.Queries
}

func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlcgen.New(pool),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, area *AffectedArea) error {
	result, err := r.queries.CreateAffectedArea(ctx, sqlcgen.CreateAffectedAreaParams{
		Name:        area.Name,
		Description: nullableString(area.Description),
		DisasterID:  area.DisasterID,
		Location:    nullableString(area.Location),
		Column5:     area.Geometry,
		Severity:    sqlcgen.AffectedAreaSeverity(area.Severity),
	})
	if err != nil {
		return database.TranslateError(err)
	}

	area.ID = strconv.FormatInt(result.ID, 10)
	area.Name = result.Name
	area.Description = textValue(result.Description)
	area.DisasterID = result.DisasterID
	area.Location = textValue(result.Location)
	area.Severity = string(result.Severity)
	if geometry, ok := result.Geometry.(string); ok {
		area.Geometry = geometry
	}
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

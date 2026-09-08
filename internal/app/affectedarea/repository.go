package affectedarea

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/apperror"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

var ErrDisasterNotFound = apperror.NotFound("disaster_not_found", "disaster not found")

type Repository interface {
	CreateArea(ctx context.Context, area *AffectedArea) error
	CreateReport(ctx context.Context, report *AffectedAreaReport) error
	UpdateAreaSeverity(ctx context.Context, id int64, severity string) error
	IncrementReportCount(ctx context.Context, id int64) error
}

type PostgresRepository struct {
	queries *sqlcgen.Queries
}

func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlcgen.New(pool),
	}
}

func (r *PostgresRepository) CreateArea(ctx context.Context, area *AffectedArea) error {
	result, err := r.queries.CreateAffectedArea(ctx, sqlcgen.CreateAffectedAreaParams{
		Name:               area.Name,
		Description:        database.NullableString(area.Description),
		DisasterID:         area.DisasterID,
		Location:           database.NullableString(area.Location),
		Column5:            area.Geometry,
		Latitude:           area.Latitude,
		Longitude:          area.Longitude,
		Severity:           sqlcgen.AffectedAreaSeverity(area.Severity),
		VerificationStatus: sqlcgen.AffectedAreaVerificationStatus(area.VerificationStatus),
	})
	if err != nil {
		return database.TranslateError(err)
	}

	area.ID = result.ID
	area.Name = result.Name
	area.Description = database.TextValue(result.Description)
	area.DisasterID = result.DisasterID
	area.Location = database.TextValue(result.Location)
	area.Geometry = geometryValue(result.Geometry)
	area.Latitude = result.Latitude
	area.Longitude = result.Longitude
	area.Severity = string(result.Severity)
	area.VerificationStatus = string(result.VerificationStatus)
	area.ReportCount = int64(result.ReportCount)
	return nil
}

func (r *PostgresRepository) CreateReport(ctx context.Context, report *AffectedAreaReport) error {
	result, err := r.queries.CreateAffectedAreaReport(ctx, sqlcgen.CreateAffectedAreaReportParams{
		AreaID:          report.AreaID,
		Name:            report.Name,
		Description:     database.NullableString(report.Description),
		LocationPayload: report.LocationPayload,
		Severity:        sqlcgen.AffectedAreaSeverity(report.Severity),
		ReporterName:    database.NullableString(report.ReporterName),
		ReporterMobile:  database.NullableString(report.ReporterMobile),
	})
	if err != nil {
		return database.TranslateError(err)
	}

	report.ID = result.ID
	report.AreaID = result.AreaID
	report.Name = result.Name
	report.Description = database.TextValue(result.Description)
	report.LocationPayload = result.LocationPayload
	report.Severity = string(result.Severity)
	report.ReporterName = database.TextValue(result.ReporterName)
	report.ReporterMobile = database.TextValue(result.ReporterMobile)
	report.CreatedAt = result.CreatedAt
	return nil
}

func (r *PostgresRepository) UpdateAreaSeverity(ctx context.Context, id int64, severity string) error {
	_, err := r.queries.UpdateAreaSeverity(ctx, sqlcgen.UpdateAreaSeverityParams{
		ID:      id,
		Column2: sqlcgen.AffectedAreaSeverity(severity),
	})
	return database.TranslateError(err)
}

func (r *PostgresRepository) IncrementReportCount(ctx context.Context, id int64) error {
	_, err := r.queries.IncrementReportCount(ctx, id)
	return database.TranslateError(err)
}

func geometryValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}

package sos

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

type Repository interface {
	Create(ctx context.Context, sos *SOS) error
	FindNearby(ctx context.Context, disasterID int32, mobile, geometry string, radiusMeters float64) (*SOS, error)
	IncrementReportCount(ctx context.Context, id int64) (*SOS, error)
}

type PostgresRepository struct {
	queries *sqlcgen.Queries
}

func NewRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		queries: sqlcgen.New(pool),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, sos *SOS) error {
	result, err := r.queries.CreateSOSRequest(ctx, sqlcgen.CreateSOSRequestParams{
		DisasterID:     sos.DisasterID,
		Column2:        pointGeoJSON(sos.Longitude, sos.Latitude),
		ReporterMobile: database.NullableString(sos.ReporterMobile),
		Message:        database.NullableString(sos.Message),
	})
	if err != nil {
		return database.TranslateError(err)
	}

	*sos = *sosFromRow(result.ID, result.DisasterID, result.Geometry, result.ReporterMobile, result.Message, result.Status, result.ReportCount, result.CreatedAt, result.UpdatedAt)
	return nil
}

func (r *PostgresRepository) FindNearby(ctx context.Context, disasterID int32, mobile, geometry string, radiusMeters float64) (*SOS, error) {
	result, err := r.queries.FindSOSNearby(ctx, sqlcgen.FindSOSNearbyParams{
		DisasterID: disasterID,
		Column2:    mobile,
		Column3:    geometry,
		Column4:    radiusMeters,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, database.TranslateError(err)
	}

	return sosFromRow(result.ID, result.DisasterID, result.Geometry, result.ReporterMobile, result.Message, result.Status, result.ReportCount, result.CreatedAt, result.UpdatedAt), nil
}

func (r *PostgresRepository) IncrementReportCount(ctx context.Context, id int64) (*SOS, error) {
	result, err := r.queries.IncrementSOSReportCount(ctx, id)
	if err != nil {
		return nil, database.TranslateError(err)
	}

	return sosFromRow(result.ID, result.DisasterID, result.Geometry, result.ReporterMobile, result.Message, result.Status, result.ReportCount, result.CreatedAt, result.UpdatedAt), nil
}

func sosFromRow(id int64, disasterID int32, geometry interface{}, reporterMobile, message *string, status sqlcgen.SosStatus, reportCount int32, createdAt, updatedAt time.Time) *SOS {
	s := &SOS{
		ID:             id,
		DisasterID:     disasterID,
		ReporterMobile: database.TextValue(reporterMobile),
		Message:        database.TextValue(message),
		Status:         string(status),
		ReportCount:    reportCount,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}
	s.Latitude, s.Longitude = coordinatesFromGeometry(geometry)
	return s
}

// pointGeoJSON builds the GeoJSON Point the database stores in its geom
// column. Coordinates are ordered as [longitude, latitude] per the GeoJSON
// spec. Returns an empty string when no location is provided, which inserts
// a NULL geom.
func pointGeoJSON(longitude, latitude *float64) string {
	if longitude == nil || latitude == nil {
		return ""
	}

	payload, err := json.Marshal(struct {
		Type        string     `json:"type"`
		Coordinates [2]float64 `json:"coordinates"`
	}{
		Type:        "Point",
		Coordinates: [2]float64{*longitude, *latitude},
	})
	if err != nil {
		return ""
	}

	return string(payload)
}

// coordinatesFromGeometry reads the latitude/longitude back out of the
// database's geom column, which PostGIS returns as GeoJSON text. A NULL geom
// yields nil values.
func coordinatesFromGeometry(value interface{}) (*float64, *float64) {
	geometry := database.GeometryToText(value)
	if geometry == "" {
		return nil, nil
	}

	var payload struct {
		Type        string     `json:"type"`
		Coordinates [2]float64 `json:"coordinates"`
	}
	if err := json.Unmarshal([]byte(geometry), &payload); err != nil {
		return nil, nil
	}
	if payload.Type != "Point" {
		return nil, nil
	}

	latitude := payload.Coordinates[1]
	longitude := payload.Coordinates[0]
	return &latitude, &longitude
}

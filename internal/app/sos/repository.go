package sos

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database"
	"github.com/techcollectiveassam/flood-response-api/internal/pkg/database/sqlcgen"
)

type Repository interface {
	Create(ctx context.Context, sos *SOS) error
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
		Column1:        pointGeoJSON(sos.Longitude, sos.Latitude),
		ReporterMobile: database.NullableString(sos.ReporterMobile),
		Message:        database.NullableString(sos.Message),
	})
	if err != nil {
		return database.TranslateError(err)
	}

	sos.ID = result.ID
	sos.Latitude, sos.Longitude = coordinatesFromGeometry(result.Geometry)
	sos.ReporterMobile = database.TextValue(result.ReporterMobile)
	sos.Message = database.TextValue(result.Message)
	sos.Status = string(result.Status)
	sos.ReportCount = result.ReportCount
	sos.CreatedAt = result.CreatedAt
	sos.UpdatedAt = result.UpdatedAt

	return nil
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

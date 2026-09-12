-- name: CreateAffectedArea :one
INSERT INTO affected_areas (
    name,
    description,
    disaster_id,
    location,
    geom,
    latitude,
    longitude,
    severity,
    centroid,
    verification_status
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    ST_SetSRID(ST_GeomFromGeoJSON(NULLIF($5::text, '')), 4326),
    $6,
    $7,
    $8,
    COALESCE(
        ST_Centroid(ST_SetSRID(ST_GeomFromGeoJSON(NULLIF($5::text, '')), 4326)),
        CASE WHEN ST_MakePoint($7, $6) IS NOT NULL
             THEN ST_SetSRID(ST_MakePoint($7, $6), 4326)
        END
    ),
    $9
)
RETURNING id, name, description, disaster_id, location,
          COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry,
          latitude, longitude, severity, verification_status, report_count;

-- name: UpdateAreaSeverity :one
UPDATE affected_areas
SET severity = $2::affected_area_severity,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING severity;

-- name: IncrementReportCount :one
UPDATE affected_areas
SET report_count = report_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING report_count;

-- name: ListAffectedAreas :many
SELECT id, name, description, disaster_id, location,
       COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry,
       latitude, longitude, severity, verification_status, report_count,
       COUNT(*) OVER() AS total_count
FROM affected_areas
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

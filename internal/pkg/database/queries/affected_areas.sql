-- name: CreateAffectedArea :one
INSERT INTO affected_areas (
    name,
    description,
    disaster_id,
    location,
    geom,
    severity
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    ST_SetSRID(ST_GeomFromGeoJSON(NULLIF($5, '')), 4326),
    $6
)
RETURNING id, name, description, disaster_id, location, COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry, severity;
-- name: CreateSOSRequest :one
INSERT INTO sos_requests (
    geom,
    reporter_mobile,
    message
)
VALUES (
    ST_SetSRID(ST_GeomFromGeoJSON(NULLIF($1::text, '')), 4326),
    $2,
    $3
)
RETURNING id,
          COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry,
          reporter_mobile, message, status, report_count, created_at, updated_at;
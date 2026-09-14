-- name: CreateSOSRequest :one
INSERT INTO sos_requests (
    disaster_id,
    geom,
    reporter_mobile,
    message
)
VALUES (
    $1,
    ST_SetSRID(ST_GeomFromGeoJSON(NULLIF($2::text, '')), 4326),
    $3,
    $4
)
RETURNING id,
          disaster_id,
          COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry,
          reporter_mobile, message, status, report_count, created_at, updated_at;

-- name: FindSOSNearby :one
SELECT id,
       disaster_id,
       COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry,
       reporter_mobile, message, status, report_count, created_at, updated_at
FROM sos_requests
WHERE disaster_id = $1
  AND status IN ('reported', 'active')
  AND geom IS NOT NULL
  AND ($2::text = '' OR reporter_mobile = $2)
  AND ST_DWithin(
        geom::geography,
        ST_SetSRID(ST_GeomFromGeoJSON($3::text), 4326)::geography,
        $4::double precision
      )
ORDER BY created_at DESC, id DESC
LIMIT 1;

-- name: IncrementSOSReportCount :one
UPDATE sos_requests
SET report_count = report_count + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id,
          disaster_id,
          COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry,
          reporter_mobile, message, status, report_count, created_at, updated_at;

-- name: ListSOSRequests :many
SELECT id,
       disaster_id,
       COALESCE(ST_AsGeoJSON(geom)::text, '') AS geometry,
       reporter_mobile, message, status, report_count, created_at, updated_at,
       COUNT(*) OVER() AS total_count
FROM sos_requests
ORDER BY created_at DESC, id DESC
LIMIT $1 OFFSET $2;
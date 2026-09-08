-- name: CreateAffectedAreaReport :one
INSERT INTO affected_area_reports (
    area_id,
    name,
    description,
    location_payload,
    severity,
    reporter_name,
    reporter_mobile
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, area_id, name, description, location_payload, severity, reporter_name, reporter_mobile, created_at;

-- name: GetAffectedAreaReportByID :one
SELECT id, area_id, name, description, location_payload, severity, reporter_name, reporter_mobile, created_at
FROM affected_area_reports
WHERE id = $1;
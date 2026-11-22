-- name: CreateAffectedAreaReport :one
INSERT INTO affected_area_reports (
    area_id,
    description,
    location_payload,
    severity,
    reporter_name,
    reporter_mobile
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, ref_id, area_id, description, location_payload, severity, reporter_name, reporter_mobile, created_at;

-- name: GetAffectedAreaReportByRefID :one
SELECT id, ref_id, area_id, description, location_payload, severity, reporter_name, reporter_mobile, created_at
FROM affected_area_reports
WHERE ref_id = $1;
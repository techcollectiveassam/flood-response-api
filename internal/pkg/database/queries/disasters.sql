-- name: CreateDisaster :one
INSERT INTO disasters (
    name,
    description,
    type,
    status,
    starts_at,
    ends_at
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    COALESCE($5, CURRENT_TIMESTAMP),
    $6
)
RETURNING id, ref_id, name, description, type, status, starts_at, ends_at;

-- name: ListDisasters :many
SELECT id, ref_id, name, description, type, status, starts_at, ends_at
FROM disasters
ORDER BY starts_at DESC;

-- name: GetDisasterByRefID :one
SELECT id, ref_id, name, description, type, status, starts_at, ends_at
FROM disasters
WHERE ref_id = $1;

-- name: GetDisasterIDByRefID :one
SELECT id
FROM disasters
WHERE ref_id = $1;

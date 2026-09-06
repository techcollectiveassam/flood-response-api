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
RETURNING id, name, description, type, status, starts_at, ends_at;

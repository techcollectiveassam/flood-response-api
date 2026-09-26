-- name: CreateCommandCenter :one
INSERT INTO command_centers (
    name,
    type,
    description,
    contact_person,
    contact_mobile,
    contact_email
)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, name, type, description, contact_person, contact_mobile, contact_email, created_at, updated_at;
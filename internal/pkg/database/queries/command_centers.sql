-- name: CreateCommandCenter :one
INSERT INTO command_centers (
    disaster_id,
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
    $6,
    $7
)
RETURNING id, disaster_id, name, type, description, contact_person, contact_mobile, contact_email, created_at, updated_at;

-- name: GetCommandCenterByDisaster :one
SELECT id, disaster_id, name, type, description, contact_person, contact_mobile, contact_email, created_at, updated_at
FROM command_centers
WHERE disaster_id = $1;

-- Guards the parent resource of the nested lookup so a missing command center
-- can be reported as an unknown disaster rather than a missing command center.
-- name: DisasterExists :one
SELECT EXISTS (SELECT 1 FROM disasters WHERE id = $1);


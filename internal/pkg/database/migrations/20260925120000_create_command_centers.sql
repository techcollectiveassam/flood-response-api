-- +goose Up
CREATE TYPE command_center_type AS ENUM (
    'government',
    'ngo',
    'group',
    'other'
);

CREATE TABLE command_centers (
    id SERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    type command_center_type NOT NULL DEFAULT 'other',
    description TEXT,

    contact_person VARCHAR(255),
    contact_mobile VARCHAR(20),
    contact_email VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE IF EXISTS command_centers;
DROP TYPE IF EXISTS command_center_type;
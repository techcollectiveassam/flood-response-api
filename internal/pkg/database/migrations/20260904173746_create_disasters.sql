-- +goose Up
CREATE TYPE disaster_type AS ENUM (
    'flood',
    'earthquake',
    'landslide'
);

CREATE TYPE disaster_status AS ENUM (
    'active',
    'resolved'
);

CREATE TABLE disasters (
    id SERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    description TEXT,

    type disaster_type NOT NULL DEFAULT 'flood',

    status disaster_status NOT NULL DEFAULT 'active',

    starts_at TIMESTAMPTZ,
    ends_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
DROP TABLE IF EXISTS disasters;
DROP TYPE IF EXISTS disaster_type;
DROP TYPE IF EXISTS disaster_status;

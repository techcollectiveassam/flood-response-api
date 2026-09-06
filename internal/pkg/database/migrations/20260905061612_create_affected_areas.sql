-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TYPE affected_area_severity AS ENUM (
    'low',
    'medium',
    'high',
    'critical'
);

CREATE TYPE affected_area_verification_status AS ENUM (
    'reported',
    'investigating',
    'verified',
    'disputed',
    'resolved'
);

CREATE TABLE affected_areas (
    id BIGSERIAL PRIMARY KEY,

    name VARCHAR(255) NOT NULL,
    description TEXT,

    disaster_id INTEGER NOT NULL REFERENCES disasters(id) ON DELETE CASCADE,

    location VARCHAR(500),
    geom geometry(Geometry, 4326),
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,

    severity affected_area_severity NOT NULL DEFAULT 'low',
    verification_status affected_area_verification_status NOT NULL DEFAULT 'reported',

    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    verified_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP

);

CREATE INDEX affected_areas_geom_idx ON affected_areas USING GIST (geom);

-- +goose Down
DROP TABLE IF EXISTS affected_areas;
DROP TYPE IF EXISTS affected_area_severity;
DROP TYPE IF EXISTS affected_area_verification_status;

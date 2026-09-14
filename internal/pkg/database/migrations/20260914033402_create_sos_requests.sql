-- +goose Up
SELECT 'up SQL query';

CREATE TYPE sos_status AS ENUM ('reported', 'active', 'resolved');

CREATE TABLE sos_requests (
    id BIGSERIAL PRIMARY KEY,
    geom geometry(Geometry, 4326),
    reporter_mobile VARCHAR(20),
    message VARCHAR(500),
    status sos_status NOT NULL DEFAULT 'reported',
    report_count INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose Down
SELECT 'down SQL query';

DROP TABLE sos_requests;
DROP TYPE sos_status;
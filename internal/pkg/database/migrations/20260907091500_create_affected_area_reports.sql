-- +goose Up
CREATE TABLE affected_area_reports (
    id               BIGSERIAL PRIMARY KEY,
    area_id          BIGINT NOT NULL REFERENCES affected_areas(id) ON DELETE CASCADE,
    name             VARCHAR(255) NOT NULL,
    description      TEXT,
    location_payload JSONB NOT NULL,
    severity         affected_area_severity NOT NULL DEFAULT 'low',
    reporter_name    VARCHAR(255),
    reporter_mobile  VARCHAR(50),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX affected_area_reports_area_id_idx ON affected_area_reports (area_id);

-- +goose Down
DROP TABLE IF EXISTS affected_area_reports;
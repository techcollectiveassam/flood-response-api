-- +goose Up
SELECT 'up SQL query';

ALTER TYPE sos_status ADD VALUE 'cancelled';

ALTER TABLE sos_requests
    ADD COLUMN disaster_id INTEGER NOT NULL
        CONSTRAINT sos_requests_disaster_id_fkey
        REFERENCES disasters(id) ON DELETE CASCADE;

CREATE INDEX sos_requests_disaster_id_idx ON sos_requests (disaster_id);
CREATE INDEX sos_requests_geom_geography_idx ON sos_requests USING GIST ((geom::geography));
-- +goose Down
SELECT 'down SQL query';

DROP INDEX IF EXISTS sos_requests_geom_geography_idx;
DROP INDEX IF EXISTS sos_requests_disaster_id_idx;

ALTER TABLE sos_requests
    DROP COLUMN disaster_id;
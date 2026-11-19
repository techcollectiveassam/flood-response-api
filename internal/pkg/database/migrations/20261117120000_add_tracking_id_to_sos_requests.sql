-- +goose Up
SELECT 'up SQL query';

ALTER TABLE sos_requests
    ADD COLUMN tracking_id UUID NOT NULL DEFAULT gen_random_uuid();

CREATE UNIQUE INDEX sos_requests_tracking_id_key ON sos_requests (tracking_id);
-- +goose Down
SELECT 'down SQL query';

DROP INDEX IF EXISTS sos_requests_tracking_id_key;

ALTER TABLE sos_requests
    DROP COLUMN tracking_id;

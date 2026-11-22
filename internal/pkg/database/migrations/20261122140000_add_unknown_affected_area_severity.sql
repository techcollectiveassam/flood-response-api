-- +goose Up
ALTER TYPE affected_area_severity ADD VALUE IF NOT EXISTS 'unknown';

-- +goose Down
-- PostgreSQL does not support removing enum values.
SELECT 'down SQL query';

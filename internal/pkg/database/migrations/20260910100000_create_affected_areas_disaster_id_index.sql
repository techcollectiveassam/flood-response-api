-- +goose Up
CREATE INDEX affected_areas_disaster_id_idx ON affected_areas (disaster_id);

-- +goose Down
DROP INDEX IF EXISTS affected_areas_disaster_id_idx;
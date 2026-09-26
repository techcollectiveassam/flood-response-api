-- +goose Up
SELECT 'up SQL query';

ALTER TABLE command_centers
    ADD COLUMN disaster_id INTEGER NOT NULL
        CONSTRAINT command_centers_disaster_id_fkey
        REFERENCES disasters(id) ON DELETE RESTRICT;

ALTER TABLE command_centers
    ADD CONSTRAINT command_centers_disaster_id_uniq UNIQUE (disaster_id);
-- +goose Down
SELECT 'down SQL query';

ALTER TABLE command_centers
    DROP CONSTRAINT IF EXISTS command_centers_disaster_id_uniq;

ALTER TABLE command_centers
    DROP COLUMN disaster_id;
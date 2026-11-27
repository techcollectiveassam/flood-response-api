-- +goose Up
ALTER TABLE affected_area_reports
    DROP COLUMN name;

ALTER TABLE affected_areas
    DROP COLUMN name,
    ADD COLUMN ref_id UUID NOT NULL DEFAULT gen_random_uuid();

CREATE UNIQUE INDEX affected_areas_ref_id_key ON affected_areas (ref_id);

-- +goose Down
DROP INDEX IF EXISTS affected_areas_ref_id_key;

ALTER TABLE affected_areas
    DROP COLUMN ref_id,
    ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';

ALTER TABLE affected_areas
    ALTER COLUMN name DROP DEFAULT;

ALTER TABLE affected_area_reports
    ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '';

ALTER TABLE affected_area_reports
    ALTER COLUMN name DROP DEFAULT;

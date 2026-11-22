-- +goose Up
ALTER TABLE affected_area_reports
    ADD COLUMN ref_id UUID NOT NULL DEFAULT gen_random_uuid();

CREATE UNIQUE INDEX affected_area_reports_ref_id_key ON affected_area_reports (ref_id);

-- +goose Down
DROP INDEX IF EXISTS affected_area_reports_ref_id_key;

ALTER TABLE affected_area_reports
    DROP COLUMN ref_id;

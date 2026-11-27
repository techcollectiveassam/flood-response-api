-- +goose Up
ALTER TABLE disasters
    ADD COLUMN ref_id UUID NOT NULL DEFAULT gen_random_uuid();

CREATE UNIQUE INDEX disasters_ref_id_key ON disasters (ref_id);

-- +goose Down
DROP INDEX IF EXISTS disasters_ref_id_key;

ALTER TABLE disasters
    DROP COLUMN ref_id;

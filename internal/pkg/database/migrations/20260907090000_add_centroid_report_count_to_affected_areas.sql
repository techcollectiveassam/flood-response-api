-- +goose Up
ALTER TABLE affected_areas
    ADD COLUMN centroid GEOMETRY(Point, 4326),
    ADD COLUMN report_count INTEGER NOT NULL DEFAULT 0;

CREATE INDEX affected_areas_centroid_idx ON affected_areas USING GIST (centroid);

-- +goose Down
DROP INDEX IF EXISTS affected_areas_centroid_idx;

ALTER TABLE affected_areas
    DROP COLUMN centroid,
    DROP COLUMN report_count;
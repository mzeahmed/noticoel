-- +goose Up
ALTER TABLE events
    RENAME COLUMN status TO severity;
ALTER TABLE events
    RENAME COLUMN data TO metadata;
ALTER TABLE events
    ADD COLUMN category TEXT;

-- +goose Down
ALTER TABLE events
    DROP COLUMN category;
ALTER TABLE events
    RENAME COLUMN metadata TO data;
ALTER TABLE events
    RENAME COLUMN severity TO status;
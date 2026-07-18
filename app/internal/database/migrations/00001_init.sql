-- +goose Up
-- Intentionally empty. Its only purpose is to bootstrap goose's own
-- version-tracking table; real schema (events, deliveries) lands in a
-- later migration once the event model exists.
SELECT 1;

-- +goose Down
SELECT 1;
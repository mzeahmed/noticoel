-- name: CreateEvent :one
INSERT INTO events (source, category, type, severity, title, message, metadata)
VALUES (?, ?, ?, ?, ?, ?, ?) RETURNING id, source, category, type, severity, title, message, metadata, created_at;

-- name: GetEvents :many
SELECT *
FROM events
ORDER BY id DESC
LIMIT ? OFFSET ?;
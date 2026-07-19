-- name: CreateEvent :one
INSERT INTO events (source, type, status, title, message, data)
VALUES (?, ?, ?, ?, ?, ?) RETURNING id, SOURCE, TYPE, status, title, message, DATA, created_at;

-- name: GetEvents :many
SELECT *
FROM events
ORDER BY id DESC
LIMIT ? OFFSET ?;
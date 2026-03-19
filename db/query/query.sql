-- name: InsertEvent :exec
INSERT INTO events (
    delivery_id,
    event_name,
    target_url,
    body_json,
    secret, 
    created_at
) VALUES (?, ?, ?, ?, ?, ?);

-- name: ListEvents :many
SELECT
    delivery_id,
    event_name,
    target_url,
    body_json,
    secret, 
    created_at
FROM events
ORDER BY created_at DESC;

-- name: GetEvent :one
SELECT
    delivery_id,
    event_name,
    target_url,
    body_json,
    secret, 
    created_at
FROM events
WHERE delivery_id = ?;
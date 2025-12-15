-- name: CreateOutboxEvent :one
INSERT INTO outbox_events (
    id, topic, event_type, event_version, key, payload
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetPendingOutboxEvents :many
SELECT * FROM outbox_events
WHERE status = 'pending' AND next_retry_at <= (now() AT TIME ZONE 'UTC')
ORDER BY created_at ASC
LIMIT $1
FOR UPDATE SKIP LOCKED;

-- name: MarkOutboxEventsSent :exec
UPDATE outbox_events
SET status = 'sent',
    sent_at = $2
WHERE id = ANY($1::uuid[]);


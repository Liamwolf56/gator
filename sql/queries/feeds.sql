-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT
    feeds.id,
    feeds.created_at,
    feeds.updated_at,
    feeds.name,
    feeds.url,
    users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT id, created_at, updated_at, name, url, user_id FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

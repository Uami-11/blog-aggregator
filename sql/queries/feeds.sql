-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, created_at, updated_at, user_id)
    VALUES ($1, $2, $3, $4, $5, $6)
RETURNING
    *;

-- name: GetFeeds :many
SELECT name, url, id FROM feeds;

-- name: GetFeedUser :one
SELECT u.name FROM feeds f
JOIN users u ON u.id = f.user_id
WHERE f.id = $1;

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

-- name: CreateFeedFollow :many
WITH inserted_feed_follows AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING *
)
SELECT inserted_feed_follows.*, f.name AS feed_name, u.name AS user_name FROM inserted_feed_follows 
JOIN users u ON u.id = inserted_feed_follows.user_id
JOIN feeds f ON f.id = inserted_feed_follows.feed_id;

-- name: UnfollowFeed :exec

DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;

-- name: FindFeedURL :one
SELECT id FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT f.name AS feed_name, u.name AS user_name FROM feed_follows ff
JOIN users u ON u.id = ff.user_id
JOIN feeds f ON f.id = ff.feed_id
WHERE u.id = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $1
WHERE id = $2;

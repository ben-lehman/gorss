-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES(
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetFeedUserName :one
SELECT users.name FROM feeds
JOIN users
ON feeds.user_id = users.id
WHERE feeds.id = $1;

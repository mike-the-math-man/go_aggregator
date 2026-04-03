-- name: GetPostsForUser :many
SELECT
    posts.*,
    feeds.name AS feed_name
FROM
    feed_follows
JOIN
    posts ON posts.feed_id = feed_follows.feed_id
JOIN
    feeds ON feeds.id = feed_follows.feed_id
WHERE
    feed_follows.user_id = $1
ORDER BY
    posts.updated_at DESC
LIMIT
    $2
;  




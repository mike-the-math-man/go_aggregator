-- name: GetFeedFollowsForUser :many
SELECT
    ff.id,
    ff.created_at,
    ff.updated_at,
    ff.user_id,
    u.name AS user_name,
    ff.feed_id,
    f.name AS feed_name
FROM
    feed_follows ff
INNER JOIN
    users u ON u.id = ff.user_id
INNER JOIN
    feeds f ON f.id = ff.feed_id
WHERE
    u.name = $1
;



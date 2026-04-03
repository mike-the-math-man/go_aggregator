-- name: GetNextFeedToFetch :one
SELECT 
    feeds.url,
    feeds.id
FROM
    feed_follows 
JOIN
    feeds ON feeds.id = feed_follows.feed_id
ORDER BY
    last_fetched_at ASC NULLS FIRST
LIMIT 1
;
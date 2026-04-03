-- name: MarkFeedFetched :exec
UPDATE 
    feed_follows 
SET 
    updated_at = NOW(), 
    last_fetched_at = NOW()
WHERE
    feed_id = $1
;
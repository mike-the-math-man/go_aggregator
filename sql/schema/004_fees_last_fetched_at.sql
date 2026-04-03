-- +goose Up
ALTER TABLE feed_follows 
ADD last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feed_follows 
DROP last_fetched_at;

-- +goose Up
CREATE TABLE IF NOT EXISTS message_views (
    message_id INT REFERENCES messages(id),
    user_id INT REFERENCES users(id),
    viewed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (message_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS message_views;

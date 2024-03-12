-- +goose Up
CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chatroom_id INT NOT NULL,
    sender_user_id INT NOT NULL,
    text TEXT NOT NULL,
    attachment_url TEXT,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    viewed BOOLEAN DEFAULT FALSE,
    deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (chatroom_id) REFERENCES chatrooms(id),
    FOREIGN KEY (sender_user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE IF EXISTS messages;

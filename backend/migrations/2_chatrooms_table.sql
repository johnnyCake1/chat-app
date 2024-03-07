-- +goose Up
CREATE TABLE IF NOT EXISTS chatrooms (
     id SERIAL PRIMARY KEY,
     is_group BOOLEAN NOT NULL,
     group_name VARCHAR(255),
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS chatrooms;

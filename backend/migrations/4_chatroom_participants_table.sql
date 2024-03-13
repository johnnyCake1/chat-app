-- +goose Up
CREATE TABLE IF NOT EXISTS chatroom_participants (
     chatroom_id INT NOT NULL,
     user_id INT NOT NULL,
     FOREIGN KEY (chatroom_id) REFERENCES chatrooms(id),
     FOREIGN KEY (user_id) REFERENCES users(id),
     PRIMARY KEY (chatroom_id, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS chatroom_participants;

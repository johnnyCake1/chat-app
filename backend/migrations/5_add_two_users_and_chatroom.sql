-- +goose Up
-- Create users
INSERT INTO users (nickname, email, password_hash, avatar_url) VALUES
    ('User1', 'user1@example.com', 'password_hash_user1', 'avatar_url_user1'),
    ('User2', 'user2@example.com', 'password_hash_user2', 'avatar_url_user2');

-- Create chatroom
INSERT INTO chatrooms (is_group, group_name) VALUES (FALSE, NULL);

-- Assign users to the chatroom
INSERT INTO chatroom_participants (chatroom_id, user_id) VALUES
    ((SELECT id FROM chatrooms ORDER BY id DESC LIMIT 1), 1),
    ((SELECT id FROM chatrooms ORDER BY id DESC LIMIT 1), 2);

-- +goose Down
-- Delete chatroom participants
DELETE FROM chatroom_participants WHERE chatroom_id = (SELECT id FROM chatrooms ORDER BY id DESC LIMIT 1);

-- Delete chatrooms
DELETE FROM chatrooms WHERE is_group = FALSE;

-- Delete users
DELETE FROM users WHERE nickname IN ('User1', 'User2');

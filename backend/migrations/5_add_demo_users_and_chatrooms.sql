-- +goose Up
CREATE TABLE IF NOT EXISTS users (
     id SERIAL PRIMARY KEY,
     nickname VARCHAR(255) NOT NULL,
     email VARCHAR(255) UNIQUE NOT NULL,
     password_hash VARCHAR(255) NOT NULL,
     avatar_url TEXT,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS chatrooms (
     id SERIAL PRIMARY KEY,
     is_group BOOLEAN NOT NULL,
     group_name VARCHAR(255),
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

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

CREATE TABLE IF NOT EXISTS chatroom_participants (
     chatroom_id INT NOT NULL,
     user_id INT NOT NULL,
     FOREIGN KEY (chatroom_id) REFERENCES chatrooms(id),
     FOREIGN KEY (user_id) REFERENCES users(id),
     PRIMARY KEY (chatroom_id, user_id)
);

-- Insert users. Hashed passwords' value is 'user123'
INSERT INTO users (nickname, email, password_hash, avatar_url, created_at) VALUES
    ('user1', 'user1@example.com', '$2a$10$eTl74BAMFMtOIVXPecQt8.FCYb2E5uU0W1aziaFG9p2wP4jhSdNh.', 'avatar1.jpg', NOW()),
    ('user2', 'user2@example.com', '$2a$10$eTl74BAMFMtOIVXPecQt8.FCYb2E5uU0W1aziaFG9p2wP4jhSdNh.', 'avatar2.jpg', NOW()),
    ('user3', 'user3@example.com', '$2a$10$eTl74BAMFMtOIVXPecQt8.FCYb2E5uU0W1aziaFG9p2wP4jhSdNh.', 'avatar3.jpg', NOW()),
    ('user4', 'user4@example.com', '$2a$10$eTl74BAMFMtOIVXPecQt8.FCYb2E5uU0W1aziaFG9p2wP4jhSdNh.', 'avatar4.jpg', NOW()),
    ('user5', 'user5@example.com', '$2a$10$eTl74BAMFMtOIVXPecQt8.FCYb2E5uU0W1aziaFG9p2wP4jhSdNh.', 'avatar5.jpg', NOW());

-- Insert chatrooms
INSERT INTO chatrooms (is_group, group_name, created_at) VALUES
    (FALSE, NULL, NOW()),
    (FALSE, NULL, NOW()),
    (FALSE, NULL, NOW());

-- Insert participants for each chatroom
INSERT INTO chatroom_participants (chatroom_id, user_id) VALUES
    (1, 1), (1, 2), -- Conversation between user1 and user2
    (2, 1), (2, 3), -- Conversation between user1 and user3
    (3, 2), (3, 4); -- Conversation between user2 and user4

-- Insert messages for the conversation between user1 and user2
INSERT INTO messages (chatroom_id, sender_user_id, text, timestamp) VALUES
    (1, 1, 'Hey user2, how are you?', NOW()),
    (1, 2, 'Hi user1, I''m good. How about you?', NOW() + interval '1 minute'),
    (1, 1, 'I''m doing well too, thanks!', NOW() + interval '2 minutes'),
    (1, 2, 'That''s great to hear!', NOW() + interval '3 minutes');

-- Insert messages for the conversation between user1 and user3
INSERT INTO messages (chatroom_id, sender_user_id, text, timestamp) VALUES
    (2, 1, 'Hello user3, how are you?', NOW()),
    (2, 3, 'Hi user1, I''m fine, thank you.', NOW() + interval '1 minute'),
    (2, 1, 'Glad to hear that!', NOW() + interval '2 minutes'),
    (2, 3, 'Yes, everything is going well.', NOW() + interval '3 minutes');

-- Insert messages for the conversation between user2 and user4
INSERT INTO messages (chatroom_id, sender_user_id, text, timestamp) VALUES
    (3, 2, 'Hey user4, how''s it going?', NOW()),
    (3, 4, 'Hi user2, I''m doing fine, thanks for asking.', NOW() + interval '1 minute'),
    (3, 2, 'That''s good to hear!', NOW() + interval '2 minutes'),
    (3, 4, 'Yes, everything is good on my end.', NOW() + interval '3 minutes');

-- +goose Up
CREATE TABLE users
(
    UserID       SERIAL PRIMARY KEY,
    Nickname     VARCHAR(255)        NOT NULL,
    Email        VARCHAR(255) UNIQUE NOT NULL,
    PasswordHash VARCHAR(255)        NOT NULL,
    AvatarURL    TEXT
);

-- +goose Down
DROP TABLE users;

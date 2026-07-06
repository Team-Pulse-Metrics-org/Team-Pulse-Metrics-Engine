-- +goose Up

ALTER TABLE users
DROP COLUMN password_hash;

ALTER TABLE users
ADD COLUMN github_id VARCHAR(255) UNIQUE,
ADD COLUMN github_username VARCHAR(255),
ADD COLUMN github_token TEXT;

-- +goose Down

ALTER TABLE users
DROP COLUMN github_id,
DROP COLUMN github_username,
DROP COLUMN github_token;

ALTER TABLE users
ADD COLUMN password_hash VARCHAR(255) NOT NULL;
-- +goose Up
ALTER TABLE users
DROP COLUMN IF EXISTS github_token;

-- +goose Down
ALTER TABLE users
ADD COLUMN github_token TEXT;

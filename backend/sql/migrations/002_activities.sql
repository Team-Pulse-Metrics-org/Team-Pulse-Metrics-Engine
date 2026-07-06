-- +goose Up

CREATE TYPE activity_type AS ENUM (
    'git_commit',
    'pull_request_closed',
    'task_completed',
    'blocker_raised'
);

CREATE TABLE activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type activity_type NOT NULL,
    payload JSONB NOT NULL,
    weight INTEGER NOT NULL DEFAULT 1,
    logged_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_activities_user_logged
ON activities(user_id, logged_at);

CREATE INDEX idx_activities_type
ON activities(type);

-- +goose Down

DROP INDEX IF EXISTS idx_activities_type;
DROP INDEX IF EXISTS idx_activities_user_logged;
DROP TABLE IF EXISTS activities;
DROP TYPE IF EXISTS activity_type;
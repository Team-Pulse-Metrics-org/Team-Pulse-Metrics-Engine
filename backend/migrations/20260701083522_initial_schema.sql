-- +goose Up

CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE user_role AS ENUM ('developer', 'lead', 'administrator');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role user_role NOT NULL DEFAULT 'developer',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_users_email ON users(email);

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

CREATE INDEX idx_activities_user_logged ON activities(user_id, logged_at);
CREATE INDEX idx_activities_type ON activities(type);

CREATE TABLE metrics_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    window_end TIMESTAMP WITH TIME ZONE NOT NULL,
    velocity_score NUMERIC(5,2) NOT NULL,
    total_commits INT NOT NULL DEFAULT 0,
    tasks_resolved INT NOT NULL DEFAULT 0,
    blockers_count INT NOT NULL DEFAULT 0,
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_user_snapshot_window
ON metrics_snapshots(user_id, window_start, window_end);

-- +goose Down

DROP TABLE IF EXISTS metrics_snapshots;
DROP TABLE IF EXISTS activities;
DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS activity_type;
DROP TYPE IF EXISTS user_role;

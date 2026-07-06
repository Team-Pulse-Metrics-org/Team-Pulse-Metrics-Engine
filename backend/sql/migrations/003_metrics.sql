-- +goose Up

CREATE TABLE metrics_snapshots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL,
    window_end TIMESTAMP WITH TIME ZONE NOT NULL,
    velocity_score NUMERIC(5, 2) NOT NULL,
    total_commits INT NOT NULL DEFAULT 0,
    tasks_resolved INT NOT NULL DEFAULT 0,
    blockers_count INT NOT NULL DEFAULT 0,
    generated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE UNIQUE INDEX idx_user_snapshot_window
ON metrics_snapshots(user_id, window_start, window_end);

-- +goose Down

DROP INDEX IF EXISTS idx_user_snapshot_window;
DROP TABLE IF EXISTS metrics_snapshots;
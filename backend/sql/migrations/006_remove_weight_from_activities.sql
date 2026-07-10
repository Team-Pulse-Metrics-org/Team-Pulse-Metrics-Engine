-- +goose Up

-- Rename blockers_count column
ALTER TABLE metrics_snapshots
RENAME COLUMN blockers_count TO open_issues;

-- Remove weight column from activities
ALTER TABLE activities
DROP COLUMN weight;

-- Rename enum value in activity type
ALTER TYPE activity_type
RENAME VALUE 'blocker_raised' TO 'open_issue';


-- +goose Down

-- Revert enum value change
ALTER TYPE activity_type
RENAME VALUE 'open_issue' TO 'blocker_raised';

-- Add weight column back
ALTER TABLE activities
ADD COLUMN weight INTEGER;

-- Revert column name change
ALTER TABLE metrics_snapshots
RENAME COLUMN open_issues TO blockers_count;
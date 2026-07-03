package models

import (
	"time"

	"github.com/google/uuid"
)

type ActivityType string

const (
	ActivityGitCommit         ActivityType = "git_commit"
	ActivityPullRequestClosed ActivityType = "pull_request_closed"
	ActivityTaskCompleted     ActivityType = "task_completed"
	ActivityBlockerRaised     ActivityType = "blocker_raised"
)

type Activities struct {
	ID         uuid.UUID    `json:"id"`
	UserID     uuid.UUID    `json:"user_id"`
	Type       ActivityType `json:"type"`
	Payload    any          `json:"payload"`
	Weight     int          `json:"weight"`
	Logged_at  time.Time    `json:"logged_at"`
	Created_at time.Time    `json:"created_at"`
}

package models

import (
	"encoding/json"
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
	ID      uuid.UUID       `json:"id"`
	UserID  uuid.UUID       `json:"user_id"`
	Type    ActivityType    `json:"type"`
	Payload json.RawMessage `json:"payload"`

	LoggedAt  time.Time `json:"logged_at"`
	CreatedAt time.Time `json:"created_at"`
}

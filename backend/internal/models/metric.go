package models

import (
	"time"

	"github.com/google/uuid"
)

type MetricsSnapshot struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	WindowStart   time.Time `json:"window_start"`
	WindowEnd     time.Time `json:"window_end"`
	VelocityScore float64   `json:"velocity_score"`
	TotalCommits  int       `json:"total_commits"`
	TasksResolved int       `json:"tasks_resolved"`
	BlockersCount int       `json:"blockers_count"`
	GeneratedAt   time.Time `json:"generated_at"`
}

type RawActivityRow struct {
	UserID   uuid.UUID
	Type     string
	LoggedAT time.Time
}

type GroupKey struct {
	UserID    uuid.UUID
	WeekStart time.Time
}

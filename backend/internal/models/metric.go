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
	OpenIssues    int       `json:"open_issues"`
	GeneratedAt   time.Time `json:"generated_at"`
}

type MetricCoordinate struct {
	Label string  `json:"label"`
	Value float64 `json:"value"`
}

type ChartTimeline struct {
	Weekly  []MetricCoordinate `json:"weekly"`
	Monthly []MetricCoordinate `json:"monthly"`
}

type UnifiedTeamMetricsResponse struct {
	Commits       ChartTimeline `json:"commits"`
	VelocityScore ChartTimeline `json:"velocity_score"`
	TasksResolved ChartTimeline `json:"tasks_resolved"`
	OpenIssues    ChartTimeline `json:"open_issues"`
}

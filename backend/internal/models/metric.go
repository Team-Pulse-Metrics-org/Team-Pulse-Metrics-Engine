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

type UnifiedMetricsResponse struct {
	Commits       ChartTimeline `json:"commits"`
	VelocityScore ChartTimeline `json:"velocity_score"`
	TasksResolved ChartTimeline `json:"tasks_resolved"`
	OpenIssues    ChartTimeline `json:"open_issues"`
}

type UserDropDownItem struct {
	UserID uuid.UUID `json:"user_id"`
	Label  string    `json:"label"`
}

func CreateUnifiedResponse(weekly []MetricsSnapshot, monthly []MetricsSnapshot) UnifiedMetricsResponse {
	commitsWeekly := make([]MetricCoordinate, 0, len(weekly))
	velocityWeekly := make([]MetricCoordinate, 0, len(weekly))
	tasksWeekly := make([]MetricCoordinate, 0, len(weekly))
	issuesWeekly := make([]MetricCoordinate, 0, len(weekly))

	commitsMonthly := make([]MetricCoordinate, 0, len(monthly))
	velocityMonthly := make([]MetricCoordinate, 0, len(monthly))
	tasksMonthly := make([]MetricCoordinate, 0, len(monthly))
	issuesMonthly := make([]MetricCoordinate, 0, len(monthly))

	for _, record := range weekly {
		weekLabel := record.WindowStart.Format("Jan _2")

		commitsWeekly = append(commitsWeekly, MetricCoordinate{Label: weekLabel, Value: float64(record.TotalCommits)})
		velocityWeekly = append(velocityWeekly, MetricCoordinate{Label: weekLabel, Value: float64(record.VelocityScore)})
		tasksWeekly = append(tasksWeekly, MetricCoordinate{Label: weekLabel, Value: float64(record.TasksResolved)})
		issuesWeekly = append(issuesWeekly, MetricCoordinate{Label: weekLabel, Value: float64(record.OpenIssues)})
	}

	for _, record := range monthly {
		monthLabel := record.WindowStart.Format("Jan")

		commitsMonthly = append(commitsMonthly, MetricCoordinate{Label: monthLabel, Value: float64(record.TotalCommits)})
		velocityMonthly = append(velocityMonthly, MetricCoordinate{Label: monthLabel, Value: float64(record.VelocityScore)})
		tasksMonthly = append(tasksMonthly, MetricCoordinate{Label: monthLabel, Value: float64(record.TasksResolved)})
		issuesMonthly = append(issuesMonthly, MetricCoordinate{Label: monthLabel, Value: float64(record.OpenIssues)})
	}

	return UnifiedMetricsResponse{
		Commits: ChartTimeline{
			Weekly:  commitsWeekly,
			Monthly: commitsMonthly,
		},
		VelocityScore: ChartTimeline{
			Weekly:  velocityWeekly,
			Monthly: velocityMonthly,
		},
		TasksResolved: ChartTimeline{
			Weekly:  tasksWeekly,
			Monthly: tasksMonthly,
		},
		OpenIssues: ChartTimeline{
			Weekly:  issuesWeekly,
			Monthly: issuesMonthly,
		},
	}

}

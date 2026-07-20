package handlers

import (
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

func HandleMetrics(c *gin.Context) {
	weeklyRecords, err := queries.GetTeamWeeklyMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "failed to collect weekly team metric snapshot"})
	}

	monthlyRecords, err := queries.GetTeamMonthlyMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to collect monthly team data metric"})
		return
	}

	var commitsWeekly, velocityWeekly, tasksWeekly, issuesWeekly []models.MetricCoordinate
	var commitsMonthly, velocityMonthly, tasksMonthly, issuesMonthly []models.MetricCoordinate

	for _, record := range weeklyRecords {
		weekLabel := record.WindowStart.Format("Jan _2")

		commitsWeekly = append(commitsWeekly, models.MetricCoordinate{Label: weekLabel, Value: float64(record.TotalCommits)})
		velocityWeekly = append(velocityWeekly, models.MetricCoordinate{Label: weekLabel, Value: float64(record.VelocityScore)})
		tasksWeekly = append(tasksWeekly, models.MetricCoordinate{Label: weekLabel, Value: float64(record.TasksResolved)})
		issuesWeekly = append(issuesWeekly, models.MetricCoordinate{Label: weekLabel, Value: float64(record.OpenIssues)})
	}

	for _, record := range monthlyRecords {
		monthLabel := record.WindowStart.Format("Jan")

		commitsMonthly = append(commitsMonthly, models.MetricCoordinate{Label: monthLabel, Value: float64(record.TotalCommits)})
		velocityMonthly = append(velocityMonthly, models.MetricCoordinate{Label: monthLabel, Value: float64(record.VelocityScore)})
		tasksMonthly = append(tasksMonthly, models.MetricCoordinate{Label: monthLabel, Value: float64(record.TasksResolved)})
		issuesMonthly = append(issuesMonthly, models.MetricCoordinate{Label: monthLabel, Value: float64(record.OpenIssues)})
	}

	responsePayload := models.UnifiedTeamMetricsResponse{
		Commits: models.ChartTimeline{
			Weekly:  commitsWeekly,
			Monthly: commitsMonthly,
		},
		VelocityScore: models.ChartTimeline{
			Weekly:  velocityWeekly,
			Monthly: velocityMonthly,
		},
		TasksResolved: models.ChartTimeline{
			Weekly:  tasksWeekly,
			Monthly: tasksMonthly,
		},
		OpenIssues: models.ChartTimeline{
			Weekly:  issuesWeekly,
			Monthly: issuesMonthly,
		},
	}
	c.JSON(http.StatusOK, responsePayload)
}

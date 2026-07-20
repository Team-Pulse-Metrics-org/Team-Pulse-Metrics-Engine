package handlers

import (
	"math"
	"net/http"
	"os"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	// Sync issues from GitHub to activities table in the background/inline
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")
	token := os.Getenv("GITHUB_PAT")
	if owner != "" && repo != "" && token != "" {
		_ = queries.SyncGithubIssues(owner, repo, token)
	}

	// 1. Fetch Stats & Activity Breakdown counts
	counts, err := queries.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch dashboard stats",
		})
		return
	}

	totalCommits := counts["git_commit"]
	prsClosed := counts["pull_request_closed"]
	tasksResolved := counts["task_completed"]
	activeBlockers := counts["open_issue"]

	// Calculate velocity score
	// Formula: velocity = [(total_commits * 1) + (tasks_completed * 5)] / (open_issues + 1)
	// Scaled to 0-100 via: 100 * (1 - e^(-raw_velocity / 59))
	var velocityScore int
	numerator := float64(totalCommits*1 + tasksResolved*5)
	denominator := float64(activeBlockers + 1)
	rawVelocity := numerator / denominator

	if rawVelocity > 0 {
		velocityScore = int(math.Round(100.0 * (1.0 - math.Exp(-rawVelocity/59.0))))
	} else {
		velocityScore = 0
	}

	start:=time.Now()

	// 2. Fetch Commit Trend
	trend, err := queries.GetCommitTrend()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch commit trend",
		})
		return
	}

	l:=middleware.LogGet()
	l.Info().Msg(time.Since(start).String())

	// 3. Fetch Top Contributors
	topContributors, err := queries.GetTopContributors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch top contributors",
		})
		return
	}

	// 4. Fetch Recent Activity
	recentActivity, err := queries.GetRecentActivity()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch recent activities",
		})
		return
	}

	// Assemble response
	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"total_commits":   totalCommits,
			"velocity_score":  velocityScore,
			"tasks_resolved":  tasksResolved,
			"active_blockers": activeBlockers,
		},
		"commit_trend": trend,
		"activity_breakdown": gin.H{
			"git_commits":          totalCommits,
			"pull_requests_closed": prsClosed,
			"tasks_resolved":      tasksResolved,
			"active_blockers":      activeBlockers,
		},
		"top_contributors": topContributors,
		"recent_activity":  recentActivity,
	})
}

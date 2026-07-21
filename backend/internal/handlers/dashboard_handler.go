package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

func (h *MetricsHandler) GetDashboard(c *gin.Context) {
	totalStart := time.Now()

	var (
		counts          map[string]int
		trend           []queries.CommitTrendItem
		topContributors []queries.TopContributorItem
		recentActivity  []queries.RecentActivityItem

		durStats           time.Duration
		durTrend           time.Duration
		durTopContributors time.Duration
		durRecentActivity  time.Duration
	)

	g, _ := errgroup.WithContext(c.Request.Context())

	// 1. Fetch Stats & Activity Breakdown counts
	g.Go(func() error {
		start := time.Now()
		var err error
		counts, err = h.q.GetDashboardStats()
		durStats = time.Since(start)
		return err
	})

	// 2. Fetch Commit Trend
	g.Go(func() error {
		start := time.Now()
		var err error
		trend, err = h.q.GetCommitTrend()
		durTrend = time.Since(start)
		return err
	})

	// 3. Fetch Top Contributors
	g.Go(func() error {
		start := time.Now()
		var err error
		topContributors, err = h.q.GetTopContributors()
		durTopContributors = time.Since(start)
		return err
	})

	// 4. Fetch Recent Activity
	g.Go(func() error {
		start := time.Now()
		var err error
		recentActivity, err = h.q.GetRecentActivity()
		durRecentActivity = time.Since(start)
		return err
	})

	if err := g.Wait(); err != nil {
		h.log.Error().Err(err).Msg("failed to fetch dashboard data from database")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch dashboard data",
		})
		return
	}

	totalDuration := time.Since(totalStart)

	h.log.Info().
		Dur("total_duration", totalDuration).
		Dur("dur_stats", durStats).
		Dur("dur_trend", durTrend).
		Dur("dur_top_contributors", durTopContributors).
		Dur("dur_recent_activity", durRecentActivity).
		Msg("dashboard metric aggregation metrics")

	totalCommits := counts["git_commit"]
	prsClosed := counts["pull_request_closed"]
	tasksResolved := counts["task_completed"]
	activeBlockers := counts["open_issue"]

	// Calculate velocity score
	var velocityScore int
	numerator := float64(totalCommits*1 + tasksResolved*5)
	denominator := float64(activeBlockers + 1)
	rawVelocity := numerator / denominator

	if rawVelocity > 0 {
		velocityScore = int(math.Round(100.0 * (1.0 - math.Exp(-rawVelocity/59.0))))
	} else {
		velocityScore = 0
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
			"tasks_resolved":       tasksResolved,
			"active_blockers":      activeBlockers,
		},
		"top_contributors": topContributors,
		"recent_activity":  recentActivity,
	})
}

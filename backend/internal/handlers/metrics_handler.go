package handlers

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type MetricsHandler struct {
	q   *queries.Queries
	cfg *config.Config
	log zerolog.Logger
}

func NewMetricsHandler(q *queries.Queries, cfg *config.Config, log zerolog.Logger) *MetricsHandler {
	return &MetricsHandler{
		q:   q,
		cfg: cfg,
		log: log,
	}
}

func (h *MetricsHandler) HandleMetrics(c *gin.Context) {
	start := time.Now()

	var wg sync.WaitGroup
	var (
		weeklyErr      error
		monthlyErr     error
		weeklyRecords  []models.MetricsSnapshot
		monthlyRecords []models.MetricsSnapshot
	)

	wg.Go(func() {
		wStart := time.Now()
		weeklyRecords, weeklyErr = h.q.GetTeamWeeklyMetrics()
		h.log.Debug().
			Dur("duration_ms", time.Since(wStart)).
			AnErr("error", weeklyErr).
			Msg("executed GetTeamWeeklyMetrics query")
	})

	wg.Go(func() {
		mStart := time.Now()
		monthlyRecords, monthlyErr = h.q.GetTeamMonthlyMetrics()
		h.log.Debug().
			Dur("duration_ms", time.Since(mStart)).
			AnErr("error", monthlyErr).
			Msg("executed GetTeamMonthlyMetrics query")
	})

	wg.Wait()
	totalDur := time.Since(start)

	if weeklyErr != nil || monthlyErr != nil {
		h.log.Error().
			Dur("total_duration_ms", totalDur).
			Err(weeklyErr).
			AnErr("monthly_err", monthlyErr).
			Msg("failed to process metrics request")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to collect metrics"})
		return
	}

	responsePayload := models.CreateUnifiedResponse(weeklyRecords, monthlyRecords)
	c.JSON(http.StatusOK, responsePayload)
}

func (h *MetricsHandler) HandleMetricDropDown(c *gin.Context) {
	users, err := h.q.GetUsersFromMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pull matching user profiles"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func (h *MetricsHandler) HandleUserMetrics(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.Error(errors.New("Path param 'id' is required"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Path param 'id' is required"})
		return
	}

	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user UUID format"})
		return
	}

	weeklyMetrics, err := h.q.GetWeeklySnapshotsByUserID(userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load historical chart data"})
		return
	}

	monthlyMetrics, err := h.q.GetMonthlySnapshotsByUserID(userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load historical chart data"})
		return
	}

	responsePayload := models.CreateUnifiedResponse(weeklyMetrics, monthlyMetrics)

	c.JSON(http.StatusOK, responsePayload)
}

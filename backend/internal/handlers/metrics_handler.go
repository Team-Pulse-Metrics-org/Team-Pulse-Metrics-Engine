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

	role := c.GetString("role")
	userIDString := c.GetString("user_id")

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	var wg sync.WaitGroup
	var (
		weeklyErr      error
		monthlyErr     error
		weeklyRecords  []models.MetricsSnapshot
		monthlyRecords []models.MetricsSnapshot
	)

	if role == "developer" {
		// Developer: Fetch user-specific metrics concurrently
		wg.Add(2)

		go func() {
			defer wg.Done()
			wStart := time.Now()
			weeklyRecords, weeklyErr = h.q.GetWeeklySnapshotsByUserID(userID)
			h.log.Debug().
				Dur("duration_ms", time.Since(wStart)).
				AnErr("error", weeklyErr).
				Msg("executed GetWeeklySnapshotsByUserID query")
		}()

		go func() {
			defer wg.Done()
			mStart := time.Now()
			monthlyRecords, monthlyErr = h.q.GetMonthlySnapshotsByUserID(userID)
			h.log.Debug().
				Dur("duration_ms", time.Since(mStart)).
				AnErr("error", monthlyErr).
				Msg("executed GetMonthlySnapshotsByUserID query")
		}()

	} else {
		// Lead / Admin: Fetch team metrics concurrently
		wg.Add(2)

		go func() {
			defer wg.Done()
			wStart := time.Now()
			weeklyRecords, weeklyErr = h.q.GetTeamWeeklyMetrics()
			h.log.Debug().
				Dur("duration_ms", time.Since(wStart)).
				AnErr("error", weeklyErr).
				Msg("executed GetTeamWeeklyMetrics query")
		}()

		go func() {
			defer wg.Done()
			mStart := time.Now()
			monthlyRecords, monthlyErr = h.q.GetTeamMonthlyMetrics()
			h.log.Debug().
				Dur("duration_ms", time.Since(mStart)).
				AnErr("error", monthlyErr).
				Msg("executed GetTeamMonthlyMetrics query")
		}()
	}

	wg.Wait()
	totalDur := time.Since(start)

	if weeklyErr != nil || monthlyErr != nil {
		h.log.Error().
			Dur("total_duration_ms", totalDur).
			Err(weeklyErr).
			AnErr("monthly_err", monthlyErr).
			Msg("failed to process metrics request")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect metrics"})
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
		c.Error(errors.New("path param 'id' is required"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "path param 'id' is required"})
		return
	}

	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user UUID format"})
		return
	}

	role := c.GetString("role")
	loggedInUser := c.GetString("user_id")

	if role == "developer" && loggedInUser != userID.String() {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access denied",
		})
		return
	}

	var wg sync.WaitGroup
	var (
		weeklyMetrics  []models.MetricsSnapshot
		monthlyMetrics []models.MetricsSnapshot
		weeklyErr      error
		monthlyErr     error
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		weeklyMetrics, weeklyErr = h.q.GetWeeklySnapshotsByUserID(userID)
	}()

	go func() {
		defer wg.Done()
		monthlyMetrics, monthlyErr = h.q.GetMonthlySnapshotsByUserID(userID)
	}()

	wg.Wait()

	if weeklyErr != nil || monthlyErr != nil {
		if weeklyErr != nil {
			c.Error(weeklyErr)
		}
		if monthlyErr != nil {
			c.Error(monthlyErr)
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load historical chart data"})
		return
	}

	responsePayload := models.CreateUnifiedResponse(weeklyMetrics, monthlyMetrics)
	c.JSON(http.StatusOK, responsePayload)
}

package handlers

import (
	"errors"
	"net/http"

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
	weeklyRecords, err := h.q.GetTeamWeeklyMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to collect weekly team metric snapshot"})
		return
	}

	monthlyRecords, err := h.q.GetTeamMonthlyMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to collect monthly team data metric"})
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

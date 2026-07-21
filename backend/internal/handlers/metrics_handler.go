package handlers

import (
	"errors"
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleMetrics(c *gin.Context) {
	weeklyRecords, err := queries.GetTeamWeeklyMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error:": "failed to collect weekly team metric snapshot"})
		return
	}

	monthlyRecords, err := queries.GetTeamMonthlyMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to collect monthly team data metric"})
		return
	}
	responsePayload := models.CreateUnifiedResponse(weeklyRecords, monthlyRecords)

	c.JSON(http.StatusOK, responsePayload)
}

func HandleMetricDropDown(c *gin.Context) {
	users, err := queries.GetUsersFromMetrics()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to pull matching user profiles"})
		return
	}

	c.JSON(http.StatusOK, users)
}

func HandleUserMetrics(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.Error(errors.New("Query param 'id' is required"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query param 'id' is required"})
		return
	}

	userID, err := uuid.Parse(idParam)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user UUID format"})
		return
	}

	weeklyMetrics, err := queries.GetWeeklySnapshotsByUserID(userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load historical chart data"})
		return
	}

	monthlyMetrics, err := queries.GetMonthlySnapshotsByUserID(userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load historical chart data"})
		return
	}

	responsePayload := models.CreateUnifiedResponse(weeklyMetrics, monthlyMetrics)

	c.JSON(http.StatusOK, responsePayload)
}

package handlers

import (
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

func GetActivities(c *gin.Context) {
	activities, err := queries.GetActivities()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch activities",
		})
		return
	}

	c.JSON(http.StatusOK, activities)
}

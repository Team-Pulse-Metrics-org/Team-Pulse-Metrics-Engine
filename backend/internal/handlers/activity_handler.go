package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *MetricsHandler) GetActivities(c *gin.Context) {
	activities, err := h.q.GetActivities()

	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch activities",
		})
		return
	}

	c.JSON(http.StatusOK, activities)
}

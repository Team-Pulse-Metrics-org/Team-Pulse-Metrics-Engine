package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetTeams(c *gin.Context) {

	developers, err := h.q.GetTeams()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to get team members."})
	}

	c.JSON(http.StatusOK, developers)
}

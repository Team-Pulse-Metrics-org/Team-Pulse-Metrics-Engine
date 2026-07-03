package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	dbqueries "github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
)

func GetAllUsers(c *gin.Context) {
	users, err := dbqueries.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := dbqueries.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

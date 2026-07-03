package handlers

import (
	"fmt"
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/gin-gonic/gin"
)

func HandlePush(c *gin.Context) {
	var payload models.PushPayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("Error binding json:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid push payload"})
		return
	}

	fmt.Printf("\nProject name: %s\n", payload.Repository.Name)
	fmt.Printf("Branch:			%s\n", payload.Ref)
	fmt.Printf("Developer 		%s\n", payload.Pusher.Name)
	fmt.Printf("Email: 		%s\n", payload.Pusher.Email)
	fmt.Printf("Commit message: %s\n", payload.Commits)

	c.JSON(http.StatusAccepted, gin.H{"status": "accepted"})
}

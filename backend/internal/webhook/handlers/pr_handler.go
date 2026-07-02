package handlers

import (
	"fmt"
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/webhook/models"
	"github.com/gin-gonic/gin"
)

func HandlePullRequest(c *gin.Context) {
	var payload models.PullRequestPayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("Error binding with json:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid PR payload"})
		return
	}

	fmt.Printf("\nPR action: 		%s\n", payload.Action)
	fmt.Printf("\nPR title:			%s\n", payload.PullRequest.Title)
	fmt.Printf("\nPR created by: 	%s\n", payload.PullRequest.User.Login)
	fmt.Printf("\nPR created at: 	%s\n", payload.PullRequest.CreatedAt)
	fmt.Printf("\nPR number 		%d\n", payload.Number)
	fmt.Printf("\nPR state 			%s\n", payload.PullRequest.State)
}

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {
	fmt.Println("Webhook received from GitHub")

	eventType := c.GetHeader("X-GitHub-Event")
	fmt.Println("Event type:", eventType)
	switch eventType {
	case "pull_request":
		HandlePullRequest(c)
	case "push":
		HandlePush(c)
	case "issues":
		HandleIssueRequest(c)
	case "ping":
		fmt.Println("pong")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		return
	}
}

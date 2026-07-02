package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleWebhook(c *gin.Context) {
	eventType := c.GetHeader("X-GitHub-Event")

	switch eventType {
	case "push":
		HandlePush(c)
	case "ping":
		fmt.Println("pong")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		return
	}
}

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

func HandlePush(c *gin.Context) {
	var payload models.PushPayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("Error binding json:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid push payload"})
		return
	}

	actor, err := queries.GetUserByGithubUsername(payload.Pusher.Name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Github account not linked",
		})
		return
	}

	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")
	commits := make([]map[string]any, 0, len(payload.Commits))

	for _, commit := range payload.Commits {
		commits = append(commits, map[string]any{
			"sha":       commit.ID,
			"message":   commit.Message,
			"timestamp": commit.Timestamp,
		})
	}

	activityPayload := map[string]any{
		"repository":   payload.Repository.Name,
		"branch":       branch,
		"author":       actor.ID,
		"commit_count": len(payload.Commits),
		"commits":      commits,
	}

	payloadJSON, err := json.Marshal(activityPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to encode payload",
		})
		return
	}

	loggedAt := time.Now()

	if len(payload.Commits) > 0 {
		loggedAt = payload.Commits[len(payload.Commits)-1].Timestamp
	}

	activity := models.Activities{
		UserID:   actor.ID,
		Type:     models.ActivityGitCommit,
		Payload:  payloadJSON,
		Weight:   len(payload.Commits),
		LoggedAt: loggedAt,
	}

	err = queries.CreateActivity(activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save activity",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Push activity stored successfully"})
}

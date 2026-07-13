package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

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


	for _, commit := range payload.Commits {

		if commit.Author.Name != payload.Pusher.Name {
			continue
		}

		activityPayload := map[string]any{
			"repository": payload.Repository.Name,
			"author":     actor.ID,
			"sha":        commit.ID,
			"message":    commit.Message,
			"timestamp":  commit.Timestamp,
		}

		payloadJSON, err := json.Marshal(activityPayload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to encode payload",
			})
			return
		}

		activity := models.Activities{
			UserID:   actor.ID,
			Type:     models.ActivityGitCommit,
			Payload:  payloadJSON,
			LoggedAt: commit.Timestamp,
		}

		err = queries.CreateActivity(activity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save activity",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Push activity stored successfully"})
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *WebhookHandler) HandlePush(c *gin.Context) {
	var payload models.PushPayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid push payload"})
		return
	}

	for _, commit := range payload.Commits {
		var actor *models.Users
		var err error

		// 1. Try to find the user in our database who is the author of this commit by their GitHub username
		if commit.Author.Username != "" {
			actor, err = h.q.GetUserByGithubUsername(commit.Author.Username)
		}

		// 2. Fallback: Try to find by email if username lookup failed or was empty
		if (err != nil || actor == nil) && commit.Author.Email != "" {
			actor, err = h.q.GetUserByEmail(commit.Author.Email)
		}

		// If the commit author is not registered in our database, skip logging this commit
		if err != nil || actor == nil {
			continue
		}

		// Prevent duplicate activities for the same commit
		existing, err := h.q.FindCommitActivityBySHA(commit.ID)
		if err == nil && existing != nil {
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
			c.Error(err)
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

		err = h.q.CreateActivity(activity)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to save activity",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Push activity stored successfully"})
}

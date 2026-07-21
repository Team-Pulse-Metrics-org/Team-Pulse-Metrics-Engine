package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *WebhookHandler) HandlePullRequest(c *gin.Context) {
	var payload models.PullRequestPayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid PR payload"})
		return
	}

	creator, err := h.q.GetUserByGithubUsername(payload.PullRequest.User.Login)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "PR creator GitHub account not linked",
		})
		return
	}

	user, err := h.q.GetUserByGithubUsername(payload.Sender.Login)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Github account not linked",
		})
		return
	}

	activityPayload := map[string]any{
		"repository":         payload.Repository.FullName,
		"pr_number":          payload.Number,
		"title":              payload.PullRequest.Title,
		"state":              payload.PullRequest.State,
		"merged":             payload.PullRequest.Merged,
		"created_by_user_id": creator.ID,
		"action_by_user_id":  user.ID,
		"source_branch":      payload.PullRequest.Head.Ref,
		"target_branch":      payload.PullRequest.Base.Ref,

		"created_at": payload.PullRequest.CreatedAt,
		"closed_at":  payload.PullRequest.ClosedAt,
		"merged_at":  payload.PullRequest.MergedAt,

		"url": payload.PullRequest.HTMLURL,
	}
	payloadJSON, err := json.Marshal(activityPayload)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to encode payload",
		})
		return
	}

	if payload.Action != "closed" {
		c.JSON(http.StatusOK, gin.H{
			"message": "PR action ignored",
		})
		return
	}

	// Check if this PR closed activity has already been logged to prevent duplicates
	existing, err := h.q.FindPRClosedActivity(payload.Number, payload.Repository.Name, payload.Repository.FullName)
	if err == nil && existing != nil {
		c.Error(err)
		c.JSON(http.StatusOK, gin.H{"message": "PR activity already stored"})
		return
	}

	loggedAt := time.Now()

	if payload.Action == "closed" {
		loggedAt = payload.PullRequest.ClosedAt
	}

	activity := models.Activities{
		UserID:  user.ID,
		Type:    models.ActivityPullRequestClosed,
		Payload: payloadJSON,

		LoggedAt: loggedAt,
	}

	err = h.q.CreateActivity(activity)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save activity",
		})
		return
	}

	h.log.Info().
		Int("pr_number", payload.Number).
		Str("repo", payload.Repository.FullName).
		Msg("PR activity stored successfully")

	c.JSON(http.StatusOK, gin.H{"message": "PR activity stored successfully"})
}

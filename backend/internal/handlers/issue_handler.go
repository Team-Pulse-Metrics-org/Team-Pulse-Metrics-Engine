package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/gin-gonic/gin"
)

func (h *WebhookHandler) HandleIssueRequest(c *gin.Context) {
	var payload models.IssuePayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid issue payload"})
		return
	}

	actor, err := h.q.GetUserByGithubUsername(payload.Sender.Login)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Github account not linked",
		})
		return
	}

	creator, err := h.q.GetUserByGithubUsername(payload.Issue.User.Login)
	var creatorID *string
	var creatorUsername string = payload.Issue.User.Login
	if err == nil && creator != nil {
		idStr := creator.ID.String()
		creatorID = &idStr
		creatorUsername = creator.GithubUsername
	}

	activityPayload := map[string]any{
		"repository":         payload.Repository.FullName,
		"issue_number":       payload.Issue.Number,
		"title":              payload.Issue.Title,
		"state":              payload.Issue.State,
		"created_by_user_id": creatorID,
		"action_by_user_id":  actor.ID,
		"created_by":         creatorUsername,
		"action_by":          actor.GithubUsername,
		"developer":          actor.GithubUsername,
		"created_at":         payload.Issue.CreatedAt,
		"closed_at":          payload.Issue.ClosedAt,
		"url":                payload.Issue.HTMLURL,
	}
	payloadJSON, err := json.Marshal(activityPayload)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to encode payload",
		})
		return
	}

	var activityType models.ActivityType
	var loggedAt time.Time

	switch payload.Action {
	case "opened", "reopened":
		activityType = models.ActivityIssueOpened
		loggedAt = payload.Issue.CreatedAt
	case "closed":
		activityType = models.ActivityTaskCompleted
		loggedAt = payload.Issue.ClosedAt

	default:
		h.log.Info().Str("action", payload.Action).Msg("Issue action ignored")
		c.JSON(http.StatusOK, gin.H{
			"message": "Issue action '" + payload.Action + "' ignored",
		})
		return
	}

	existingActivity, err := h.q.FindIssueActivity(payload.Issue.Number, payload.Repository.Name, payload.Repository.FullName)
	if err == nil && existingActivity != nil {
		existingActivity.Type = activityType
		existingActivity.Payload = payloadJSON
		existingActivity.LoggedAt = loggedAt

		err = h.q.UpdateActivity(*existingActivity)
		if err != nil {
			c.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update activity",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Issue activity updated successfully",
		})
		return
	}

	activity := models.Activities{
		UserID:   actor.ID,
		Type:     activityType,
		Payload:  payloadJSON,
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

	c.JSON(http.StatusOK, gin.H{
		"message": "Issue activity stored successfully",
	})
}

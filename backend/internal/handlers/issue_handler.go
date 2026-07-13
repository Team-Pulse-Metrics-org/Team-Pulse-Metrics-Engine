package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

func HandleIssueRequest(c *gin.Context) {
	var payload models.IssuePayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid PR payload"})
		return
	}

	actor, err := queries.GetUserByGithubUsername(payload.Sender.Login)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Github account not linked",
		})
		return
	}

	creator, err := queries.GetUserByGithubUsername(payload.Issue.User.Login)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Issue creator GitHub account not linked",
		})
		return
	}

	activityPayload := map[string]any{
		"repository":   payload.Repository.FullName,
		"issue_number": payload.Issue.Number,
		"title":        payload.Issue.Title,
		"state":        payload.Issue.State,

		"created_by_user_id": creator.ID,
		"action_by_user_id":  actor.ID,

		"created_at": payload.Issue.CreatedAt,
		"closed_at":  payload.Issue.ClosedAt,

		"url": payload.Issue.HTMLURL,
	}
	payloadJSON, err := json.Marshal(activityPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to encode payload",
		})
		return
	}

	var activityType models.ActivityType
	var loggedAt time.Time

	switch payload.Action {
	case "opened":
		activityType = models.ActivityIssueOpened
		loggedAt = payload.Issue.CreatedAt
	case "closed":
		activityType = models.ActivityTaskCompleted
		loggedAt = payload.Issue.ClosedAt

	default:
		c.JSON(http.StatusOK, gin.H{
			"message": "Issue action ignored",
		})
		return
	}

	activity := models.Activities{
		UserID:   actor.ID,
		Type:     activityType,
		Payload:  payloadJSON,
		LoggedAt: loggedAt,
	}
	err = queries.CreateActivity(activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save activity",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Issue activity stored successfully",
	})
}

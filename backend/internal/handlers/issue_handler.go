package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

func HandleIssueRequest(c *gin.Context) {
	var payload models.IssuePayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid issue payload"})
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
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Issue action '%s' ignored", payload.Action),
		})
		return
	}

	existingActivity, err := queries.FindIssueActivity(payload.Issue.Number, payload.Repository.Name, payload.Repository.FullName)
	if err == nil && existingActivity != nil {
		existingActivity.Type = activityType
		existingActivity.Payload = payloadJSON
		existingActivity.LoggedAt = loggedAt

		err = queries.UpdateActivity(*existingActivity)
		if err != nil {
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


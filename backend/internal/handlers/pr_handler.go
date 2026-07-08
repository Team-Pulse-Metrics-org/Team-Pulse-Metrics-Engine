package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

func HandlePullRequest(c *gin.Context) {
	var payload models.PullRequestPayload

	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid PR payload"})
		return
	}

	user, err := queries.GetUserByGithubUsername(payload.Sender.Login)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Github account not linked",
		})
		return
	}

	activityPayload := map[string]any{
		"repository": payload.Repository.FullName,
		"pr_number":  payload.Number,
		"title":      payload.PullRequest.Title,

		"action": payload.Action,
		"state":  payload.PullRequest.State,
		"merged": payload.PullRequest.Merged,

		"created_by": payload.PullRequest.User.Login,
		"action_by":  payload.Sender.Login,

		"source_branch": payload.PullRequest.Head.Ref,
		"target_branch": payload.PullRequest.Base.Ref,

		"created_at": payload.PullRequest.CreatedAt,
		"closed_at":  payload.PullRequest.ClosedAt,
		"merged_at":  payload.PullRequest.MergedAt,

		"url": payload.PullRequest.HTMLURL,
	}

	payloadJSON, err := json.Marshal(activityPayload)
	if err != nil {
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

	loggedAt := time.Now()

	if payload.Action == "closed" {
		loggedAt = payload.PullRequest.ClosedAt
	}

	activity := models.Activities{
		UserID:   user.ID,
		Type:     models.ActivityPullRequestClosed,
		Payload:  payloadJSON,
		Weight:   1,
		LoggedAt: loggedAt,
	}

	err = queries.CreateActivity(activity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save activity",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "PR activity stored successfully"})
}

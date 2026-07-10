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

		"developer":     user.FirstName + " " + user.LastName,
		"created_by":    user.FirstName + " " + user.LastName,
		"action_by":     user.FirstName + " " + user.LastName,
		"source_branch": payload.PullRequest.Head.Ref,
		"target_branch": payload.PullRequest.Base.Ref,

		"created_at": payload.PullRequest.CreatedAt,
		"closed_at":  payload.PullRequest.ClosedAt,
		"merged_at":  payload.PullRequest.MergedAt,

		"url": payload.PullRequest.HTMLURL,
	}
	fmt.Println("Developer being stored:", user.FirstName+" "+user.LastName)
	fmt.Println("Payload author:", activityPayload["author"])
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
		UserID:  user.ID,
		Type:    models.ActivityPullRequestClosed,
		Payload: payloadJSON,

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

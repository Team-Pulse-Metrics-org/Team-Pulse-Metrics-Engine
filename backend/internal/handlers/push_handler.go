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

	fmt.Println("========== PUSH WEBHOOK RECEIVED ==========")

	// Bind JSON
	if err := c.ShouldBindBodyWithJSON(&payload); err != nil {
		fmt.Println("❌ JSON Bind Error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid push payload"})
		return
	}

	fmt.Println("✅ JSON parsed successfully")
	fmt.Printf("Sender Login : %q\n", payload.Sender.Login)
	fmt.Printf("Pusher Name  : %q\n", payload.Pusher.Name)
	fmt.Printf("Repository   : %q\n", payload.Repository.Name)
	fmt.Printf("Branch Ref   : %q\n", payload.Ref)
	fmt.Printf("Commit Count : %d\n", len(payload.Commits))

	// Lookup user
	fmt.Println("🔍 Looking up GitHub user...")

	user, err := queries.GetUserByGithubUsername(payload.Sender.Login)
	if err != nil {
		fmt.Println("❌ User lookup failed:", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Github account not linked",
		})
		return
	}

	fmt.Printf("✅ User found: %s (%s)\n", user.ID, user.GithubUsername)

	// Clean branch name
	branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

	commits := make([]map[string]any, 0, len(payload.Commits))

	for _, commit := range payload.Commits {
		fmt.Printf("Commit -> SHA: %s | Message: %s\n", commit.ID, commit.Message)

		commits = append(commits, map[string]any{
			"sha":       commit.ID,
			"message":   commit.Message,
			"timestamp": commit.Timestamp,
		})
	}

	activityPayload := map[string]any{
		"repository":   payload.Repository.Name,
		"branch":       branch,
		"author":       payload.Pusher.Name,
		"author_email": payload.Pusher.Email,
		"commit_count": len(payload.Commits),
		"commits":      commits,
	}

	fmt.Println("📦 Marshaling activity payload...")

	payloadJSON, err := json.Marshal(activityPayload)
	if err != nil {
		fmt.Println("❌ Marshal Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to encode payload",
		})
		return
	}

	fmt.Println("✅ Payload marshaled")

	loggedAt := time.Now()

	if len(payload.Commits) > 0 {
		loggedAt = payload.Commits[len(payload.Commits)-1].Timestamp
	}

	activity := models.Activities{
		UserID:   user.ID,
		Type:     models.ActivityGitCommit,
		Payload:  payloadJSON,
		Weight:   len(payload.Commits),
		LoggedAt: loggedAt,
	}

	fmt.Println("💾 Inserting activity into database...")

	err = queries.CreateActivity(activity)
	if err != nil {
		fmt.Println("❌ Database Insert Error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println("✅ Activity inserted successfully!")
	fmt.Println("==========================================")

	c.JSON(http.StatusOK, gin.H{
		"message": "Push activity stored successfully",
	})
}

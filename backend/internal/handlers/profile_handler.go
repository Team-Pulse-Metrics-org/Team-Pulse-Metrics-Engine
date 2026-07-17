package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type ProfileResponse struct {
	GithubID  int    `json:"github_id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Followers int    `json:"followers"`
	Following int    `json:"following"`
}

func GetGitHubProfile(c *gin.Context) {
	// 1. Get logged in user ID from Gin context (set by AuthRequired middleware)
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "unauthorized"})
		return
	}

	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "invalid user session"})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid user id format"})
		return
	}

	// 2. Retrieve user from database
	user, err := queries.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "user not found"})
		return
	}

	// 3. Fetch GitHub profile details using API
	ghProfileUrl := "https://api.github.com/users/" + user.GithubUsername
	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", ghProfileUrl, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to create request to GitHub"})
		return
	}

	// Use GITHUB_TOKEN or GITHUB_PAT if available for authentication/rate limit
	tokenStr := os.Getenv("GITHUB_TOKEN")
	if tokenStr == "" {
		tokenStr = os.Getenv("GITHUB_PAT")
	}
	if tokenStr != "" {
		req.Header.Set("Authorization", "token "+tokenStr)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to fetch github profile"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "github api error"})
		return
	}

	var ghResponse struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
		Followers int    `json:"followers"`
		Following int    `json:"following"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ghResponse); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to parse github profile"})
		return
	}

	// 4. Resolve email
	email := ghResponse.Email
	if email == "" {
		if user.Email != "" {
			email = user.Email
		} else {
			// Call the fetchPrimaryEmail function from github_login_handle.go
			var client *http.Client
			if tokenStr != "" {
				config := &oauth2.Config{}
				tok := &oauth2.Token{AccessToken: tokenStr}
				client = config.Client(c.Request.Context(), tok)
			} else {
				client = http.DefaultClient
			}
			primaryEmail, err := fetchPrimaryEmail(client)
			if err != nil {
				c.Error(err)
			} else if primaryEmail != "" {
				email = primaryEmail
			}
		}
	}

	// Fallback for name if it's empty in GitHub API
	name := ghResponse.Name
	if name == "" {
		name = user.FirstName
		if user.LastName != "" {
			name += " " + user.LastName
		}
	}
	if name == "" {
		name = user.GithubUsername
	}

	// Resolve githubID
	githubID := ghResponse.ID
	if githubID == 0 && user.GithubID != "" {
		if idVal, err := strconv.Atoi(user.GithubID); err == nil {
			githubID = idVal
		}
	}

	// Prepare profile response
	profile := ProfileResponse{
		GithubID:  githubID,
		Username:  user.GithubUsername,
		Name:      name,
		Email:     email,
		AvatarURL: ghResponse.AvatarURL,
		Followers: ghResponse.Followers,
		Following: ghResponse.Following,
	}

	c.JSON(http.StatusOK, profile)
}
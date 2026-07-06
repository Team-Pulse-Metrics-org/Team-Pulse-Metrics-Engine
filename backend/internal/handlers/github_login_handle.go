package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/auth"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func HandleGithubLogin(c *gin.Context) {
	var req models.GithubLoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid request payload"})
		return
	}

	config := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Endpoint:     github.Endpoint,
	}

	token, err := config.Exchange(context.Background(), req.Code)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Failed to exchange Github code"})
		return
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to fetch github profile"})
		return
	}
	defer resp.Body.Close()

	var ghResponse models.GithubResponse

	json.NewDecoder(resp.Body).Decode(&ghResponse)

	//Check database if user exists or not
	var userID string
	var userEmail string
	var userRole string

	appToken, err := auth.GenerateJWTToken(userID, userRole)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to generate session token"})
		return
	}

	c.JSON(http.StatusOK, models.LoginResponse{
		Status: "success",
		Token:  appToken,
		User: models.UserDetails{
			ID:    userID,
			Email: userEmail,
			Role:  userRole,
		},
	})

}

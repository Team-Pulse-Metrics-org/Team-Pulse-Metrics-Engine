package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/auth"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
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
	user, err := queries.GetUserByEmail(ghResponse.Email)
	if err != nil {

		// User doesn't exist, create a new one
		if err == sql.ErrNoRows {

			firstName := ghResponse.Name
			lastName := ""

			// Split full name into first and last name (if available)
			nameParts := strings.Fields(ghResponse.Name)
			if len(nameParts) > 0 {
				firstName = nameParts[0]
			}
			if len(nameParts) > 1 {
				lastName = strings.Join(nameParts[1:], " ")
			}

			newUser := &models.Users{
				Email:          ghResponse.Email,
				FirstName:      firstName,
				LastName:       lastName,
				Role:           models.RoleDeveloper,
				GithubID:       strconv.Itoa(ghResponse.ID),
				GithubUsername: ghResponse.Login,
			}

			user, err = queries.CreateUser(newUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "failed to create user",
				})
				return
			}

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "database error",
			})
			return
		}
	}

	userID := user.ID.String()
	userEmail := user.Email
	userRole := string(user.Role)

	//need to write the logic for find or create user

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

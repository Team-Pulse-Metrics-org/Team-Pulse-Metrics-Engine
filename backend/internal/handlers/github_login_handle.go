package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/auth"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type AuthHandler struct {
	q   *queries.Queries
	cfg *config.Config
	log zerolog.Logger
}

func NewAuthHandler(q *queries.Queries, cfg *config.Config, log zerolog.Logger) *AuthHandler {
	return &AuthHandler{
		q:   q,
		cfg: cfg,
		log: log,
	}
}

func (h *AuthHandler) HandleGithubLogin(c *gin.Context) {
	var req models.GithubLoginRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid request payload"})
		return
	}

	oauthConfig := &oauth2.Config{
		ClientID:     h.cfg.GithubClientID,
		ClientSecret: h.cfg.GithubClientSecret,
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email"},
	}

	token, err := oauthConfig.Exchange(context.Background(), req.Code)
	if err != nil {
		h.log.Warn().Err(err).Msg("failed to exchange GitHub OAuth code")
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Failed to exchange Github code"})
		return
	}

	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to fetch github profile"})
		return
	}
	defer resp.Body.Close()

	var ghResponse models.GithubResponse

	if err := json.NewDecoder(resp.Body).Decode(&ghResponse); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to parse github profile"})
		return
	}

	if ghResponse.Email == "" {
		if email, err := fetchPrimaryEmail(c.Request.Context(), client); err != nil {
			h.log.Warn().Err(err).Msg("failed to fetch primary email from GitHub API")
		} else {
			ghResponse.Email = email
		}
	}

	//Check database if user exists or not
	githubID := strconv.Itoa(ghResponse.ID)
	user, err := h.q.GetUserByGithubID(githubID)

	if err != nil {

		// User doesn't exist, create a new one
		if errors.Is(err, sql.ErrNoRows) {
			h.log.Info().Str("github_id", githubID).Msg("registering new user via GitHub OAuth")

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
				GithubID:       githubID,
				GithubUsername: ghResponse.Login,
			}

			user, err = h.q.CreateUser(newUser)
			if err != nil {
				h.log.Error().Err(err).Msg("failed to save new OAuth user to database")
				c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "failed to create user",
				})
				return
			}

		} else {
			h.log.Error().Err(err).Str("github_id", githubID).Msg("database query failed during GitHub login")
			c.Error(err)
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

	appToken, err := auth.GenerateJWTToken(userID, userRole)
	if err != nil {
		h.log.Error().Err(err).Str("user_id", userID).Msg("failed to sign JWT session token")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "failed to generate session token"})
		return
	}

	h.log.Info().Str("user_id", userID).Msg("user authenticated successfully via GitHub")
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

func fetchPrimaryEmail(ctx context.Context, client *http.Client) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []models.GithubEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	for _, e := range emails {
		if e.Verified {
			return e.Email, nil
		}
	}

	return "", nil
}

// func HandleLogout(c *gin.Context) {

// }

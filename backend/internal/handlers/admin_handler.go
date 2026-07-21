package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
)

type AdminHandler struct {
	q   *queries.Queries
	cfg *config.Config
	log zerolog.Logger
}

func NewAdminHandler(q *queries.Queries, cfg *config.Config, log zerolog.Logger) *AdminHandler {
	return &AdminHandler{
		q:   q,
		cfg: cfg,
		log: log,
	}
}

// =====================================================
// GetUsers
// GET /api/v1/admin/users
//
// Returns all users from the database.
// Used by the Admin page to populate the users table.
// =====================================================
func (h *AdminHandler) GetUsers(c *gin.Context) {

	// Fetch all users from database
	users, err := h.q.GetAllUsers()

	// If query fails return error response
	if err != nil {
		h.log.Error().Err(err).Msg("failed to fetch all users")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch users",
		})
		return
	}

	// Return users as JSON
	c.JSON(http.StatusOK, users)
}

type UpdateRoleRequest struct {
	Role string `json:"role"`
}

func (h *AdminHandler) UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := h.q.UpdateUserRole(userID, req.Role)
	if err != nil {
		h.log.Error().Err(err).Str("user_id", userID).Str("role", req.Role).Msg("failed to update user role")
		c.Error(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	h.log.Info().Str("user_id", userID).Str("role", req.Role).Msg("user role updated successfully")
	c.JSON(200, gin.H{
		"message": "role updated successfully",
	})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	// Get user ID from URL
	id := c.Param("id")

	// Delete from database
	err := h.q.DeleteUser(id)
	if err != nil {
		h.log.Error().Err(err).Str("user_id", id).Msg("failed to delete user")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete user",
		})
		return
	}

	h.log.Info().Str("user_id", id).Msg("user deleted successfully")
	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	createdUser, err := h.q.CreateUser(&user)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to manually create user")
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	h.log.Info().Str("user_id", createdUser.ID.String()).Msg("user manually created successfully")
	c.JSON(http.StatusCreated, createdUser)
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
)

// =====================================================
// GetUsers
// GET /api/v1/admin/users
//
// Returns all users from the database.
// Used by the Admin page to populate the users table.
// =====================================================
func GetUsers(c *gin.Context) {

	// Fetch all users from database
	users, err := queries.GetAllUsers()

	// If query fails return error response
	if err != nil {
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

func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var req struct {
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	err := queries.UpdateUserRole(userID, req.Role)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"message": "role updated successfully",
	})
}
func DeleteUser(c *gin.Context) {
	// Get user ID from URL
	id := c.Param("id")

	// Delete from database
	err := queries.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user deleted successfully",
	})
}
func CreateUser(c *gin.Context) {
	var user models.Users

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	createdUser, err := queries.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

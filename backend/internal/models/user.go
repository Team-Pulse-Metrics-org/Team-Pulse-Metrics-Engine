package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleDeveloper     UserRole = "developer"
	RoleLead          UserRole = "lead"
	RoleAdministrator UserRole = "administrator"
)

type Users struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Role         UserRole  `json:"role"`

	GithubID       string    `json:"github_id"`
	GithubUsername string    `json:"github_username"`
	GithubToken    string    `json:"github_token"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

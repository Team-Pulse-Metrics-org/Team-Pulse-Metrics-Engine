package queries

import (
	"github.com/google/uuid"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
)

func GetUserByID(id uuid.UUID) (*models.Users, error) {
	var user models.Users

	query := `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			role,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func GetUserByEmail(email string) (*models.Users, error) {
	var user models.Users

	query := `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			role,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
func GetAllUsers() ([]models.Users, error) {
	query := `
		SELECT
			id,
			email,
			password_hash,
			first_name,
			last_name,
			role,
			created_at,
			updated_at
		FROM users
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.Users

	for rows.Next() {
		var user models.Users

		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

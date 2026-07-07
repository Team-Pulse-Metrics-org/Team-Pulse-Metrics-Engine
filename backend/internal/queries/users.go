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
	
			first_name,
			last_name,
			role,
			github_id,
    		github_username,
    
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,

		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.GithubID,
		&user.GithubUsername,

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
			
			first_name,
			last_name,
			role,
			github_id,
            github_username,
    
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	err := database.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,

		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.GithubID,
		&user.GithubUsername,

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
			
			first_name,
			last_name,
			role,
			github_id,
            github_username,
           
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

			&user.FirstName,
			&user.LastName,
			&user.Role,
			&user.GithubID,
			&user.GithubUsername,

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

func CreateUser(user *models.Users)(*models.Users, error){
	query := `
			INSERT INTO users (
				email,
				first_name,
				last_name,
				role,
				github_id,
				github_username
			)
			VALUES (
				$1, $2, $3, $4, $5, $6
			)
			RETURNING
				id,
				email,
				first_name,
				last_name,
				role,
				github_id,
				github_username,
				created_at,
				updated_at
			`	
	err := database.DB.QueryRow(
		query,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Role,
		user.GithubID,
		user.GithubUsername,
	).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.GithubID,
		&user.GithubUsername,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
	
}

func GetUserByGithubUsername(username string) (*models.Users,error){
	var user models.Users

	query:=`SELECT
			id,
			email,
			first_name,
			last_name,
			role,
			github_id,
			github_username,
			created_at,
			updated_at
		FROM users
		WHERE LOWER(github_username) = LOWER($1);`
	
	
	err := database.DB.QueryRow(query, username).Scan(
		&user.ID,
		&user.Email,

		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.GithubID,
		&user.GithubUsername,

		&user.CreatedAt,
		&user.UpdatedAt,
	)
		if err != nil {
		return nil, err
	}

	return &user, nil
}
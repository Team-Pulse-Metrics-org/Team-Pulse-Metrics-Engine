package queries

import (
	"github.com/google/uuid"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
)

func (q *Queries) GetUserByID(id uuid.UUID) (*models.Users, error) {
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

	err := q.db.QueryRow(query, id).Scan(
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
func (q *Queries) GetUserByEmail(email string) (*models.Users, error) {
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

	err := q.db.QueryRow(query, email).Scan(
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
func (q *Queries) GetAllUsers() ([]models.Users, error) {
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

	rows, err := q.db.Query(query)
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

func (q *Queries) CreateUser(user *models.Users) (*models.Users, error) {
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
	err := q.db.QueryRow(
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

func (q *Queries) GetUserByGithubUsername(username string) (*models.Users, error) {
	var user models.Users

	query := `SELECT
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
		WHERE github_username = $1;`

	err := q.db.QueryRow(query, username).Scan(
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

func (q *Queries) GetUserByGithubID(githubID string) (*models.Users, error) {
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
		WHERE github_id = $1
	`

	err := q.db.QueryRow(query, githubID).Scan(
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
func (q *Queries) UpdateUserRole(id string, role string) error {
	query := `
		UPDATE users
		SET role = $1
		WHERE id = $2
	`

	_, err := q.db.Exec(query, role, id)
	if err != nil {
		return err
	}

	return nil
}
func (q *Queries) DeleteUser(id string) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := q.db.Exec(query, id)
	return err
}

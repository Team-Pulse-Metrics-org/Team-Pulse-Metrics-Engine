package queries

import (
	"github.com/google/uuid"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
)

func GetActivityByID(id uuid.UUID) (*models.Activities, error) {
	var activity models.Activities

	query := `
		SELECT
			a.id,
			a.user_id,
			a.type,
			a.payload,
			a.logged_at,
			a.created_at,
			COALESCE(u.first_name || ' ' || u.last_name, u.github_username, 'Unknown') AS developer_name
		FROM activities a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Type,
		&activity.Payload,
		&activity.LoggedAt,
		&activity.CreatedAt,
		&activity.DeveloperName,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func GetActivitiesByUserID(userID uuid.UUID) ([]models.Activities, error) {
	query := `
		SELECT
			a.id,
			a.user_id,
			a.type,
			a.payload,
			a.logged_at,
			a.created_at,
			COALESCE(u.first_name || ' ' || u.last_name, u.github_username, 'Unknown') AS developer_name
		FROM activities a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.user_id = $1
		ORDER BY a.logged_at DESC
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activities

	for rows.Next() {
		var activity models.Activities

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.Type,
			&activity.Payload,
			&activity.LoggedAt,
			&activity.CreatedAt,
			&activity.DeveloperName,
		)
		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

func GetAllActivities() ([]models.Activities, error) {
	query := `
		SELECT
			a.id,
			a.user_id,
			a.type,
			a.payload,
			a.logged_at,
			a.created_at,
			COALESCE(u.first_name || ' ' || u.last_name, u.github_username, 'Unknown') AS developer_name
		FROM activities a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.logged_at DESC
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activities

	for rows.Next() {
		var activity models.Activities

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.Type,
			&activity.Payload,
			&activity.LoggedAt,
			&activity.CreatedAt,
			&activity.DeveloperName,
		)
		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

func CreateActivity(activity models.Activities) error {
	query := `INSERT INTO activities (
			user_id,
			type,
			payload,
			logged_at
		)
		VALUES ($1, $2, $3, $4)`

	_, err := database.DB.Exec(
		query,
		activity.UserID,
		activity.Type,
		activity.Payload,

		activity.LoggedAt,
	)
	return err
}
func GetActivities() ([]models.Activities, error) {
	rows, err := database.DB.Query(`
    SELECT
        a.id,
        a.user_id,
        a.type,
        a.payload,
        a.logged_at,
        a.created_at,
        u.first_name || ' ' || u.last_name AS developer_name
    FROM activities a
    LEFT JOIN users u
    ON a.user_id = u.id
    ORDER BY a.logged_at DESC
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []models.Activities

	for rows.Next() {
		var activity models.Activities

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.Type,
			&activity.Payload,
			&activity.LoggedAt,
			&activity.CreatedAt,
			&activity.DeveloperName,
		)

		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

func FindIssueActivity(issueNumber int, repoName string, repoFullName string) (*models.Activities, error) {
	var activity models.Activities

	query := `
		SELECT
			id,
			user_id,
			type,
			payload,
			logged_at,
			created_at
		FROM activities
		WHERE (type = 'open_issue' OR type = 'task_completed')
		  AND (payload->>'issue_number')::int = $1
		  AND (payload->>'repository' = $2 OR payload->>'repository' = $3)
		LIMIT 1
	`

	err := database.DB.QueryRow(query, issueNumber, repoName, repoFullName).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Type,
		&activity.Payload,
		&activity.LoggedAt,
		&activity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func UpdateActivity(activity models.Activities) error {
	query := `
		UPDATE activities
		SET type = $1,
		    payload = $2,
		    logged_at = $3
		WHERE id = $4
	`

	_, err := database.DB.Exec(
		query,
		activity.Type,
		activity.Payload,
		activity.LoggedAt,
		activity.ID,
	)
	return err
}

func FindCommitActivityBySHA(sha string) (*models.Activities, error) {
	var activity models.Activities

	query := `
		SELECT
			id,
			user_id,
			type,
			payload,
			logged_at,
			created_at
		FROM activities
		WHERE type = 'git_commit'
		  AND payload->>'sha' = $1
		LIMIT 1
	`

	err := database.DB.QueryRow(query, sha).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Type,
		&activity.Payload,
		&activity.LoggedAt,
		&activity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}

func FindPRClosedActivity(prNumber int, repoName string, repoFullName string) (*models.Activities, error) {
	var activity models.Activities

	query := `
		SELECT
			id,
			user_id,
			type,
			payload,
			logged_at,
			created_at
		FROM activities
		WHERE type = 'pull_request_closed'
		  AND (payload->>'pr_number')::int = $1
		  AND (payload->>'repository' = $2 OR payload->>'repository' = $3)
		LIMIT 1
	`

	err := database.DB.QueryRow(query, prNumber, repoName, repoFullName).Scan(
		&activity.ID,
		&activity.UserID,
		&activity.Type,
		&activity.Payload,
		&activity.LoggedAt,
		&activity.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &activity, nil
}




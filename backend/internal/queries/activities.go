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
			id,
			user_id,
			type,
			payload,
			
			logged_at,
			created_at
		FROM activities
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
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

func GetActivitiesByUserID(userID uuid.UUID) ([]models.Activities, error) {
	query := `
		SELECT
			id,
			user_id,
			type,
			payload,
		
			logged_at,
			created_at
		FROM activities
		WHERE user_id = $1
		ORDER BY logged_at DESC
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
			id,
			user_id,
			type,
			payload,
			
			logged_at,
			created_at
		FROM activities
		ORDER BY logged_at DESC
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
		VALUES ($1, $2, $3, $4, $5)`

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
        SELECT id, user_id, type, payload, logged_at, created_at
        FROM activities
        ORDER BY logged_at DESC
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
		)

		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

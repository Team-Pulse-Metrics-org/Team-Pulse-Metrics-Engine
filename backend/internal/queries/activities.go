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
			weight,
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
		&activity.Weight,
		&activity.Logged_at,
		&activity.Created_at,
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
			weight,
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
			&activity.Weight,
			&activity.Logged_at,
			&activity.Created_at,
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
			weight,
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
			&activity.Weight,
			&activity.Logged_at,
			&activity.Created_at,
		)
		if err != nil {
			return nil, err
		}

		activities = append(activities, activity)
	}

	return activities, nil
}

package queries

import (
	"github.com/google/uuid"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
)

func GetMetricByID(id uuid.UUID) (*models.MetricsSnapshot, error) {
	var metric models.MetricsSnapshot

	query := `
		SELECT
			id,
			user_id,
			window_start,
			window_end,
			velocity_score,
			total_commits,
			tasks_resolved,
			blockers_count,
			generated_at
		FROM metrics_snapshots
		WHERE id = $1
	`

	err := database.DB.QueryRow(query, id).Scan(
		&metric.ID,
		&metric.UserID,
		&metric.WindowStart,
		&metric.WindowEnd,
		&metric.VelocityScore,
		&metric.TotalCommits,
		&metric.TasksResolved,
		&metric.BlockersCount,
		&metric.GeneratedAt,
	)

	if err != nil {
		return nil, err
	}

	return &metric, nil
}
func GetMetricsByUserID(userID uuid.UUID) ([]models.MetricsSnapshot, error) {
	query := `
		SELECT
			id,
			user_id,
			window_start,
			window_end,
			velocity_score,
			total_commits,
			tasks_resolved,
			blockers_count,
			generated_at
		FROM metrics_snapshots
		WHERE user_id = $1
	`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.MetricsSnapshot

	for rows.Next() {
		var metric models.MetricsSnapshot

		err := rows.Scan(
			&metric.ID,
			&metric.UserID,
			&metric.WindowStart,
			&metric.WindowEnd,
			&metric.VelocityScore,
			&metric.TotalCommits,
			&metric.TasksResolved,
			&metric.BlockersCount,
			&metric.GeneratedAt,
		)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}
func GetAllMetrics() ([]models.MetricsSnapshot, error) {
	query := `
		SELECT
			id,
			user_id,
			window_start,
			window_end,
			velocity_score,
			total_commits,
			tasks_resolved,
			blockers_count,
			generated_at
		FROM metrics_snapshots
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []models.MetricsSnapshot

	for rows.Next() {
		var metric models.MetricsSnapshot

		err := rows.Scan(
			&metric.ID,
			&metric.UserID,
			&metric.WindowStart,
			&metric.WindowEnd,
			&metric.VelocityScore,
			&metric.TotalCommits,
			&metric.TasksResolved,
			&metric.BlockersCount,
			&metric.GeneratedAt,
		)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	return metrics, nil
}

func GetWeeklySnapshots(userID uuid.UUID) ([]models.MetricsSnapshot, error) {
	query := `	
		SELECT id, user_id, window_start, window_end, velocity_score, total_commits, tasks_resolved, blockers_count, generated_at
		FROM (
				SELECT * FROM metrics_snapshots
				WHERE user_id = $1
				ORDER BY window_start DESC
				LIMIT 7
		) sub 
		ORDER BY window_start ASC;	
		`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.MetricsSnapshot
	for rows.Next() {
		var s models.MetricsSnapshot
		err := rows.Scan(&s.ID, &s.UserID, &s.WindowStart, &s.WindowEnd, &s.VelocityScore, &s.TotalCommits, &s.TasksResolved, &s.BlockersCount, &s.GeneratedAt)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

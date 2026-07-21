package queries

import (
	"context"

	"github.com/google/uuid"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
)

func (q *Queries) GetWeeklySnapshotsByUserID(userID uuid.UUID) ([]models.MetricsSnapshot, error) {
	query := `	
		SELECT id, user_id, window_start, window_end, velocity_score, total_commits, tasks_resolved, open_issues, generated_at
		FROM (
				SELECT * FROM metrics_snapshots
				WHERE user_id = $1
				ORDER BY window_start DESC
				LIMIT 7
		) sub 
		ORDER BY window_start ASC;	
		`

	rows, err := q.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.MetricsSnapshot
	for rows.Next() {
		var s models.MetricsSnapshot
		err := rows.Scan(&s.ID, &s.UserID, &s.WindowStart, &s.WindowEnd, &s.VelocityScore, &s.TotalCommits, &s.TasksResolved, &s.OpenIssues, &s.GeneratedAt)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (q *Queries) GetMonthlySnapshotsByUserID(userID uuid.UUID) ([]models.MetricsSnapshot, error) {
	query := `
        SELECT 
            COALESCE(SUM(total_commits), 0) AS total_commits,
            COALESCE(ROUND(AVG(velocity_score):: numeric, 1), 0) AS velocity_score,
            COALESCE(SUM(tasks_resolved), 0) AS tasks_resolved,
            COALESCE(SUM(open_issues), 0) AS open_issues,
            DATE_TRUNC('month', window_start) AS targeted_month
        FROM metrics_snapshots
        WHERE user_id = $1
          AND DATE_TRUNC('month', window_start) IN (
                SELECT DISTINCT DATE_TRUNC('month', window_start)  
                FROM metrics_snapshots
                WHERE user_id = $1
                ORDER BY DATE_TRUNC('month', window_start) DESC
                LIMIT 6
        )
        GROUP BY targeted_month
        ORDER BY targeted_month ASC;
    `

	rows, err := q.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.MetricsSnapshot
	for rows.Next() {
		var s models.MetricsSnapshot
		err := rows.Scan(&s.TotalCommits, &s.VelocityScore, &s.TasksResolved, &s.OpenIssues, &s.WindowStart)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (q *Queries) CreateMetric(ctx context.Context) error {
	l := middleware.LogGet()
	query := `
	WITH weekly_records AS(
		SELECT user_id,
			DATE_TRUNC('week', logged_at) AS window_start,
			DATE_TRUNC('week', logged_at) + INTERVAL '7 days' AS window_end,
			COUNT(CASE WHEN type = 'git_commit' THEN 1 END) AS total_commits,
			COUNT(CASE WHEN type = 'open_issue' THEN 1 END) AS open_issues,
			COUNT(CASE WHEN type = 'task_completed' THEN 1 END) AS tasks_resolved
		FROM activities
		WHERE user_id = $1
		GROUP BY 
			user_id,
			DATE_TRUNC('week', logged_at)
	),
	calculated_velocity AS (
		SELECT *, (((total_commits * 1) + (tasks_resolved * 5))::float / (open_issues + 1)::float) AS raw_velocity
		FROM weekly_records
	),
	final_metrics AS (
	SELECT *,
	ROUND (LEAST(((raw_velocity / 70.0) * 100), 100.0):: numeric, 2) AS velocity_percentage
	FROM calculated_velocity
	)
	INSERT INTO metrics_snapshots
	(user_id, window_start, window_end, velocity_score, total_commits, tasks_resolved, open_issues)
	SELECT user_id, window_start, window_end, velocity_percentage, total_commits, tasks_resolved, open_issues
	FROM final_metrics
	ON CONFLICT (user_id, window_start, window_end)
	DO UPDATE SET
		open_issues = EXCLUDED.open_issues,
		tasks_resolved = EXCLUDED.tasks_resolved,
		total_commits = EXCLUDED.total_commits,
		velocity_score = EXCLUDED.velocity_score,
		generated_at = CURRENT_TIMESTAMP;
	`

	userIDs, err := q.GetAllUserIDFromActivity(ctx)
	if err != nil {
		return err
	}

	stmt, err := q.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, userID := range userIDs {
		_, err := stmt.ExecContext(ctx, userID)
		if err != nil {
			l.Warn().Timestamp().Msgf("failed running metrics for user %s: %v", userID, err)
			continue
		}
	}
	return nil
}

func (q *Queries) GetTeamWeeklyMetrics() ([]models.MetricsSnapshot, error) {
	query := `
		SELECT 
				COALESCE(SUM(total_commits), 0) AS total_commits,
				COALESCE(ROUND(AVG(velocity_score):: numeric, 1), 0) AS velocity_score,
				COALESCE(SUM(tasks_resolved), 0) AS tasks_resolved,
				COALESCE(SUM(open_issues), 0) AS open_issues,
				window_start,
				window_end
		FROM metrics_snapshots
		WHERE window_start IN (
				SELECT DISTINCT window_start
				FROM metrics_snapshots
				ORDER BY window_start DESC
				LIMIT 7
		)
		GROUP BY window_start, window_end
		ORDER BY window_start ASC;
	`

	rows, err := q.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.MetricsSnapshot
	for rows.Next() {
		var s models.MetricsSnapshot
		err := rows.Scan(&s.TotalCommits, &s.VelocityScore, &s.TasksResolved, &s.OpenIssues, &s.WindowStart, &s.WindowEnd)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (q *Queries) GetTeamMonthlyMetrics() ([]models.MetricsSnapshot, error) {
	query := `
			SELECT 
				COALESCE(SUM(total_commits), 0) AS total_commits,
				COALESCE(ROUND(AVG(velocity_score):: numeric, 1), 0) AS velocity_score,
				COALESCE(SUM(tasks_resolved), 0) AS tasks_resolved,
				COALESCE(SUM(open_issues), 0) AS open_issues,
				DATE_TRUNC('month', window_start) AS targeted_month
		FROM metrics_snapshots
		WHERE DATE_TRUNC('month', window_start) IN (
				SELECT DISTINCT DATE_TRUNC('month',window_start)	
				FROM metrics_snapshots
				ORDER BY DATE_TRUNC('month',window_start) DESC
				LIMIT 6
		)
		GROUP BY DATE_TRUNC('month', window_start)
		ORDER BY targeted_month ASC;
	`

	rows, err := q.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snapshots []models.MetricsSnapshot
	for rows.Next() {
		var s models.MetricsSnapshot
		err := rows.Scan(&s.TotalCommits, &s.VelocityScore, &s.TasksResolved, &s.OpenIssues, &s.WindowStart)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, s)
	}
	return snapshots, nil
}

func (q *Queries) GetUsersFromMetrics() ([]models.UserDropDownItem, error) {
	query := `
			SELECT DISTINCT ms.user_id, (u.first_name || ' ' || u.last_name) AS full_name
			FROM metrics_snapshots ms
			INNER JOIN users u ON ms.user_id = u.id
			ORDER BY full_name ASC;
	`

	rows, err := q.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserDropDownItem
	for rows.Next() {
		var item models.UserDropDownItem
		if err := rows.Scan(&item.UserID, &item.Label); err != nil {
			return nil, err
		}
		users = append(users, item)
	}
	return users, nil
}

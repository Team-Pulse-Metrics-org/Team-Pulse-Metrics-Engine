package queries

import (
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
)

func (q *Queries) GetTeams() ([]models.TeamMember, error) {
	query := `
SELECT
    u.id,
    CONCAT(u.first_name, ' ', u.last_name) AS name,

    COUNT(CASE WHEN a.type = 'git_commit' THEN 1 END) AS commits,

    COUNT(CASE WHEN a.type = 'task_completed' THEN 1 END) AS tasks_resolved,

    ROUND(
        (
            COUNT(CASE WHEN a.type = 'git_commit' THEN 1 END) * 1.0 +
            COUNT(CASE WHEN a.type = 'task_completed' THEN 1 END) * 2.0
        )::numeric,
        2
    ) AS velocity_score,

    0 AS open_issues

FROM users u
LEFT JOIN activities a
    ON u.id = a.user_id

GROUP BY u.id, u.first_name, u.last_name
ORDER BY velocity_score DESC;
`
	rows, err := q.db.Query(query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var developers []models.TeamMember

	for rows.Next() {
		var dev models.TeamMember

		err := rows.Scan(
			&dev.ID,
			&dev.Name,
			&dev.Commits,
			&dev.TasksResolved,
			&dev.Velocity,
			&dev.OpenIssues,
		)

		if err != nil {
			return nil, err
		}

		developers = append(developers, dev)
	}
	return developers, nil
}

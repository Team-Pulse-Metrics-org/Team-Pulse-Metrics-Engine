package handlers

import (
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/gin-gonic/gin"
)

type TeamMember struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Commits       int     `json:"commits"`
	Velocity      float64 `json:"velocity"`
	TasksResolved int     `json:"tasksResolved"`
	OpenIssues    int     `json:"openIssues"`
}

func GetTeams(c *gin.Context) {

	rows, err := database.DB.Query(`
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
`)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer rows.Close()

	var developers []TeamMember

	for rows.Next() {
		var dev TeamMember

		err := rows.Scan(
			&dev.ID,
			&dev.Name,
			&dev.Commits,
			&dev.TasksResolved,
			&dev.Velocity,
			&dev.OpenIssues,
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		developers = append(developers, dev)
	}

	c.JSON(http.StatusOK, developers)
}

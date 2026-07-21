package models

type TeamMember struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Commits       int     `json:"commits"`
	Velocity      float64 `json:"velocity"`
	TasksResolved int     `json:"tasksResolved"`
	OpenIssues    int     `json:"openIssues"`
}

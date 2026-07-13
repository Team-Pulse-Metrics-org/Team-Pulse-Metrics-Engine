package models

import "time"

type IssuePayload struct {
	Action string `json:"action"`

	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`

	Repository struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`

	Issue Issue `json:"issue"`
}

type Issue struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`

	User struct {
		Login string `json:"login"`
	} `json:"user"`

	CreatedAt time.Time `json:"created_at"`
	ClosedAt  time.Time `json:"closed_at"`
}

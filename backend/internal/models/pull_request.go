package models

import "time"

type PullRequestPayload struct {
	Action string `json:"action"`
	Number int    `json:"number"`
	Sender struct {
		Login string `json:"login"`
	} `json:"sender"`

	PullRequest PullRequest `json:"pull_request"`
	Repository  struct {
		Name     string `json:"name"`
		FullName string `json:"full_name"`
	} `json:"repository"`
}

type PullRequest struct {
	Title   string `json:"title"`
	State   string `json:"state"`
	HTMLURL string `json:"html_url"`
	Merged  bool   `json:"merged"`

	User struct {
		Login string `json:"login"`
	} `json:"user"`

	Head struct {
		Ref string `json:"ref"`
	} `json:"head"`

	Base struct {
		Ref string `json:"ref"`
	} `json:"base"`

	CreatedAt time.Time `json:"created_at"`
	ClosedAt  time.Time `json:"closed_at"`
	MergedAt  time.Time `json:"merged_at"`
}

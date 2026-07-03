package models

import "time"

type PullRequestPayload struct {
	Action      string      `json:"action"`
	Number      int         `json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository  struct {
		Name string `json:"name"`
	} `json:"repository"`
}

type PullRequest struct {
	Title string `json:"title"`
	State string `json:"state"`
	User  struct {
		Login string `json:"login"`
	} `json:"user"`
	CreatedAt time.Time `json:"created_at"`
}

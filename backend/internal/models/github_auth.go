package models

type GithubLoginRequest struct {
	Code string `json:"code" binding:"required"`
}

type LoginResponse struct {
	Status string      `json:"status"`
	Token  string      `json:"token"`
	User   UserDetails `json:"user"`
}

type UserDetails struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type GithubResponse struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

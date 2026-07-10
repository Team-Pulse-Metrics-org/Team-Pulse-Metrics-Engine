package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/joho/godotenv"
)

type CommitResponse []struct {
	SHA string `json:"sha"`

	Commit struct {
		Message string `json:"message"`

		Author struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Date  string `json:"date"`
		} `json:"author"`
	} `json:"commit"`

	Author *struct {
		Login string `json:"login"`
	} `json:"author"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	database.ConnectDB()

	token := os.Getenv("GITHUB_PAT")
	owner := os.Getenv("GITHUB_OWNER")
	repo := os.Getenv("GITHUB_REPO")

	client := &http.Client{}

	page := 1
	imported := 0
	skipped := 0

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/commits?per_page=100&page=%d",
			owner,
			repo,
			page,
		)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			panic(fmt.Sprintf("GitHub API Error: %s\n%s", resp.Status, string(body)))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			panic(err)
		}

		var commits CommitResponse

		if err := json.Unmarshal(body, &commits); err != nil {
			panic(err)
		}

		// No more commits
		if len(commits) == 0 {
			break
		}

		fmt.Printf("Importing page %d (%d commits)\n", page, len(commits))

		for _, commit := range commits {

			var actor *models.Users
			var err error

			if commit.Author != nil {
				// Normal case: GitHub could identify the user
				actor, err = queries.GetUserByGithubUsername(commit.Author.Login)
			} else {
				// Fallback: GitHub couldn't identify the user,
				// use the Git author name instead
				fmt.Printf("Fallback lookup using Git author name: %s\n", commit.Commit.Author.Name)

				actor, err = queries.GetUserByGithubUsername(commit.Commit.Author.Name)
			}

			if err != nil {
				fmt.Printf("User lookup failed: %v\n", err)
				skipped++
				continue
			}

			activityPayload := map[string]any{
				"repository": repo,
				"author":     actor.ID,
				"sha":        commit.SHA,
				"message":    commit.Commit.Message,
				"timestamp":  commit.Commit.Author.Date,
			}

			payloadJSON, err := json.Marshal(activityPayload)
			if err != nil {
				fmt.Println("JSON marshal error:", err)
				skipped++
				continue
			}

			loggedAt, err := time.Parse(time.RFC3339, commit.Commit.Author.Date)
			if err != nil {
				fmt.Println("Time parse error:", commit.Commit.Author.Date, err)
				skipped++
				continue
			}

			activity := models.Activities{
				UserID:   actor.ID,
				Type:     models.ActivityGitCommit,
				Payload:  payloadJSON,
				LoggedAt: loggedAt,
			}

			if err := queries.CreateActivity(activity); err != nil {
				fmt.Println("CreateActivity failed:", err)
				skipped++
				continue
			}

			fmt.Println("Inserted:", commit.SHA)
			imported++
		}

		page++
	}

	fmt.Println("===================================")
	fmt.Println("History Sync Completed")
	fmt.Println("Imported :", imported)
	fmt.Println("Skipped  :", skipped)
	fmt.Println("===================================")

	// PR Merged retriver

	page = 1
	imported = 0
	skipped = 0

	type PullRequestResponse []struct {
		Number int `json:"number"`
		State  string `json:"state"`
		MergedAt *time.Time `json:"merged_at"`

		User struct {
			Login string `json:"login"`
		} `json:"user"`

		Title string `json:"title"`
		HTMLURL string `json:"html_url"`

		Head struct {
			Ref string `json:"ref"`
		} `json:"head"`

		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`

		MergedBy *struct {
			Login string `json:"login"`
		} `json:"merged_by"`
		

		CreatedAt time.Time `json:"created_at"`
		ClosedAt  time.Time `json:"closed_at"`
	}

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/pulls?state=closed&per_page=100&page=%d",
			owner,
			repo,
			page,
		)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			panic(fmt.Sprintf("GitHub API Error: %s\n%s", resp.Status, string(body)))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			panic(err)
		}

		var pullrequest PullRequestResponse

		if err := json.Unmarshal(body, &pullrequest); err != nil {
			panic(err)
		}

		if len(pullrequest) == 0 {
			break
		}

		fmt.Printf("Importing PR page %d (%d PRs)\n", page, len(pullrequest))

		for _, pr := range pullrequest {

			// Skip PRs that were closed but never merged
			if pr.MergedAt == nil {
				skipped++
				continue
			}

			creator, err := queries.GetUserByGithubUsername(pr.User.Login)
			if err != nil {
				skipped++
				fmt.Printf("User lookup failed: %v\n", err)
				continue
			}


			actionUser := creator

			if pr.MergedBy != nil {
				if merger, err := queries.GetUserByGithubUsername(pr.MergedBy.Login); err == nil {
					actionUser = merger
				}
			}

			
			activityPayload := map[string]any{
				"repository":         repo,
				"pr_number":          pr.Number,
				"title":              pr.Title,
				"state":              pr.State,

				"created_by_user_id": creator.ID,
				"action_by_user_id": actionUser.ID,
				
				"source_branch": pr.Head.Ref,
				"target_branch": pr.Base.Ref,

				"created_at": pr.CreatedAt,
				"closed_at":  pr.ClosedAt,
				"merged_at":  *pr.MergedAt,

				"url": pr.HTMLURL,
			}

			payloadJSON, err := json.Marshal(activityPayload)
			if err != nil {
				skipped++
				fmt.Println(err)
				continue
			}

			activity := models.Activities{
				UserID:   actionUser.ID,
				Type:     models.ActivityPullRequestClosed,
				Payload:  payloadJSON,
				LoggedAt: *pr.MergedAt,
			}

			if err := queries.CreateActivity(activity); err != nil {
				skipped++
				fmt.Println("CreateActivity failed:", err)
				continue
			}

			imported++
			fmt.Printf("Imported PR #%d\n", pr.Number)
		}

		page++
	}
	fmt.Println("===================================")
	fmt.Println("PR History Sync Completed")
	fmt.Println("Imported :", imported)
	fmt.Println("Skipped  :", skipped)
	fmt.Println("===================================")

	// --------------------------------------------
// Issue History Sync
// --------------------------------------------

	page = 1
	imported = 0
	skipped = 0

	type IssueResponse []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
		HTMLURL string `json:"html_url"`

		User struct {
			Login string `json:"login"`
		} `json:"user"`

		PullRequest *struct {
			URL string `json:"url"`
		} `json:"pull_request"`

		CreatedAt time.Time  `json:"created_at"`
		ClosedAt  *time.Time `json:"closed_at"`
	}

	for {

		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/issues?state=all&per_page=100&page=%d",
			owner,
			repo,
			page,
		)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			panic(fmt.Sprintf("GitHub API Error: %s\n%s", resp.Status, string(body)))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			panic(err)
		}

		var issues IssueResponse

		if err := json.Unmarshal(body, &issues); err != nil {
			panic(err)
		}

		if len(issues) == 0 {
			break
		}

		fmt.Printf("Importing Issue page %d (%d issues)\n", page, len(issues))

		for _, issue := range issues {

			// Skip pull requests
			if issue.PullRequest != nil {
				continue
			}

			creator, err := queries.GetUserByGithubUsername(issue.User.Login)
			if err != nil {
				fmt.Printf("User lookup failed: %v\n", err)
				skipped++
				continue
			}

			activityPayload := map[string]any{
				"repository":         repo,
				"issue_number":       issue.Number,
				"title":              issue.Title,
				"state":              issue.State,

				"created_by_user_id": creator.ID,
				"action_by_user_id":  creator.ID,

				"created_at": issue.CreatedAt,
				"closed_at":  issue.ClosedAt,

				"url": issue.HTMLURL,
			}

			payloadJSON, err := json.Marshal(activityPayload)
			if err != nil {
				skipped++
				continue
			}

			var (
				activityType models.ActivityType
				loggedAt     time.Time
			)

			if issue.State == "open" {
				activityType = models.ActivityIssueOpened
				loggedAt = issue.CreatedAt
			} else {
				activityType = models.ActivityTaskCompleted

				if issue.ClosedAt != nil {
					loggedAt = *issue.ClosedAt
				} else {
					loggedAt = issue.CreatedAt
				}
			}

			activity := models.Activities{
				UserID:   creator.ID,
				Type:     activityType,
				Payload:  payloadJSON,
				LoggedAt: loggedAt,
			}

			if err := queries.CreateActivity(activity); err != nil {
				fmt.Println("CreateActivity failed:", err)
				skipped++
				continue
			}

			imported++
			fmt.Printf("Imported Issue #%d\n", issue.Number)
		}

		page++
	}

	fmt.Println("===================================")
	fmt.Println("Issue History Sync Completed")
	fmt.Println("Imported :", imported)
	fmt.Println("Skipped  :", skipped)
	fmt.Println("===================================")
}

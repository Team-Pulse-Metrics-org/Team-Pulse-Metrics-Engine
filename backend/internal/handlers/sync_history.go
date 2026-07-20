package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
)

type SyncCountSummary struct {
	Imported int `json:"imported"`
	Skipped  int `json:"skipped"`
}

type SyncSummary struct {
	Commits      SyncCountSummary `json:"commits"`
	PullRequests SyncCountSummary `json:"pull_requests"`
	Issues       SyncCountSummary `json:"issues"`
}

func RunSync() {
	l := middleware.LogGet()
	defer func() {
		if r := recover(); r != nil {
			l.Error().Msgf("RunSync recovered from panic: %v", r)
		}
	}()

	l.Info().Msg("Starting GitHub synchronization...")

	cfg := LoadGitHubConfig()
	commitSummary := SyncCommits(cfg)
	prSummary := SyncPR(cfg)
	issueSummary := SyncIssue(cfg)

	l.Info().
		Int("commits_imported", commitSummary.Imported).
		Int("commits_skipped", commitSummary.Skipped).
		Int("prs_imported", prSummary.Imported).
		Int("prs_skipped", prSummary.Skipped).
		Int("issues_imported", issueSummary.Imported).
		Int("issues_skipped", issueSummary.Skipped).
		Msg("GitHub synchronization completed")

	if err := UpdateLastSyncGist(time.Now()); err != nil {
		l.Error().Err(err).Msg("Failed to update last sync gist")
	}
}

func HandleSync(c *gin.Context) {
	go RunSync()

	c.JSON(http.StatusOK, gin.H{
		"message": "Sync started",
	})
}
//getting last sync function
func GetLastSync(c *gin.Context) {
    sync, err := ReadLastSyncGist()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, sync)
}

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

type GithubSyncConfig struct {
	Token  string
	Owner  string
	Repo   string
	Client *http.Client
}

var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

func LoadGitHubConfig() GithubSyncConfig {
	return GithubSyncConfig{
		Token:  os.Getenv("GITHUB_PAT"),
		Owner:  os.Getenv("GITHUB_OWNER"),
		Repo:   os.Getenv("GITHUB_REPO"),

		Client: httpClient,
	}
}

func SyncCommits(cfg GithubSyncConfig) SyncCountSummary {
	owner := cfg.Owner
	repo := cfg.Repo
	token := cfg.Token
	client := cfg.Client

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

			// Skip if commit already exists
			existing, err := queries.FindCommitActivityBySHA(commit.SHA)
			if err == nil && existing != nil {
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

	return SyncCountSummary{Imported: imported, Skipped: skipped}
}

// PR Merged retriver
func SyncPR(cfg GithubSyncConfig) SyncCountSummary {
	owner := cfg.Owner
	repo := cfg.Repo
	token := cfg.Token
	client := cfg.Client

	page := 1
	imported := 0
	skipped := 0

	type PullRequestResponse []struct {
		Number   int        `json:"number"`
		State    string     `json:"state"`
		MergedAt *time.Time `json:"merged_at"`

		User struct {
			Login string `json:"login"`
		} `json:"user"`

		Title   string `json:"title"`
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

			// Skip if PR activity already exists
			existing, err := queries.FindPRClosedActivity(pr.Number, repo, owner+"/"+repo)
			if err == nil && existing != nil {
				skipped++
				continue
			}

			activityPayload := map[string]any{
				"repository": repo,
				"pr_number":  pr.Number,
				"title":      pr.Title,
				"state":      pr.State,

				"created_by_user_id": creator.ID,
				"action_by_user_id":  actionUser.ID,

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

	return SyncCountSummary{Imported: imported, Skipped: skipped}
}

// --------------------------------------------
// Issue History Sync
// --------------------------------------------
func SyncIssue(cfg GithubSyncConfig) SyncCountSummary {
	owner := cfg.Owner
	repo := cfg.Repo
	token := cfg.Token
	client := cfg.Client

	page := 1
	imported := 0
	skipped := 0

	type IssueResponse []struct {
		Number  int    `json:"number"`
		Title   string `json:"title"`
		State   string `json:"state"`
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

			// Skip if issue activity already exists
			existing, err := queries.FindIssueActivity(issue.Number, repo, owner+"/"+repo)
			if err == nil && existing != nil {
				skipped++
				continue
			}

			activityPayload := map[string]any{
				"repository":   repo,
				"issue_number": issue.Number,
				"title":        issue.Title,
				"state":        issue.State,

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

	return SyncCountSummary{Imported: imported, Skipped: skipped}
}

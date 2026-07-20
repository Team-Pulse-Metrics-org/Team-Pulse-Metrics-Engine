package handlers

import (
	"context"
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
	"golang.org/x/sync/errgroup"
)

type SyncCountSummary struct {
	Returned int `json:"returned"`
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

	// 1. Read last sync timestamp from Gist
	var lastSyncTime time.Time
	lastSyncStr := "None (Full Sync)"
	syncObj, err := ReadLastSyncGist()
	if err != nil {
		l.Warn().Err(err).Msg("Failed to read last sync gist; performing full sync")
	} else if syncObj.LastSynced != "" {
		if t, err := time.Parse(time.RFC3339, syncObj.LastSynced); err == nil {
			lastSyncTime = t
			lastSyncStr = syncObj.LastSynced
		} else {
			l.Warn().Err(err).Msg("Failed to parse last sync timestamp; performing full sync")
		}
	}

	l.Info().Msg(fmt.Sprintf("Last Sync: %s", lastSyncStr))

	// 2. Capture sync start time
	syncStartTime := time.Now().UTC()

	cfg := LoadGitHubConfig()

	// 3. Run sync operations concurrently
	g, _ := errgroup.WithContext(context.Background())

	g.Go(func() error {
		_, err := SyncCommits(cfg, lastSyncTime)
		return err
	})

	g.Go(func() error {
		_, err := SyncPR(cfg, lastSyncTime)
		return err
	})

	g.Go(func() error {
		_, err := SyncIssue(cfg, lastSyncTime)
		return err
	})

	if err := g.Wait(); err != nil {
		l.Error().Err(err).Msg("Synchronization failed; Gist timestamp will not be updated")
		return
	}

	l.Info().Msg("Synchronization completed successfully.")

	// 4. Update Gist only after all sync operations succeed
	if err := UpdateLastSyncGist(syncStartTime); err != nil {
		l.Error().Err(err).Msg("Failed to update last sync gist after successful sync")
	} else {
		l.Info().Msg("Updated last_sync in GitHub Gist.")
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
	token := os.Getenv("GITHUB_PAT")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	return GithubSyncConfig{
		Token:  token,
		Owner:  os.Getenv("GITHUB_OWNER"),
		Repo:   os.Getenv("GITHUB_REPO"),

		Client: httpClient,
	}
}

func SyncCommits(cfg GithubSyncConfig, since time.Time) (SyncCountSummary, error) {
	l := middleware.LogGet()
	owner := cfg.Owner
	repo := cfg.Repo
	token := cfg.Token
	client := cfg.Client

	page := 1
	returnedTotal := 0
	imported := 0
	skipped := 0

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/commits?per_page=100&page=%d",
			owner,
			repo,
			page,
		)
		if !since.IsZero() {
			url += fmt.Sprintf("&since=%s", since.UTC().Format(time.RFC3339))
		}

		l.Info().Msg(fmt.Sprintf("GitHub request URL: %s", url))

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "Team-Pulse-Metrics-Engine")

		resp, err := client.Do(req)
		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, fmt.Errorf("GitHub API Error: %s\n%s", resp.Status, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		var commits CommitResponse

		if err := json.Unmarshal(body, &commits); err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		returnedTotal += len(commits)

		if len(commits) == 0 {
			break
		}

		reachedOldCommits := false
		for _, commit := range commits {
			if !since.IsZero() {
				commitDate, parseErr := time.Parse(time.RFC3339, commit.Commit.Author.Date)
				if parseErr == nil && (commitDate.Before(since) || commitDate.Equal(since)) {
					reachedOldCommits = true
					break
				}
			}

			var actor *models.Users

			if commit.Author != nil {
				actor, err = queries.GetUserByGithubUsername(commit.Author.Login)
			} else {
				actor, err = queries.GetUserByGithubUsername(commit.Commit.Author.Name)
			}

			if err != nil {
				skipped++
				continue
			}

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
				skipped++
				continue
			}

			loggedAt, err := time.Parse(time.RFC3339, commit.Commit.Author.Date)
			if err != nil {
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
				skipped++
				continue
			}

			imported++
		}

		if len(commits) < 100 || reachedOldCommits {
			break
		}

		page++
	}

	l.Info().Msg(fmt.Sprintf("Commits returned by GitHub: %d", returnedTotal))
	l.Info().Msg(fmt.Sprintf("New commits inserted: %d", imported))

	return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, nil
}

// PR Merged retriver
func SyncPR(cfg GithubSyncConfig, since time.Time) (SyncCountSummary, error) {
	l := middleware.LogGet()
	owner := cfg.Owner
	repo := cfg.Repo
	token := cfg.Token
	client := cfg.Client

	page := 1
	returnedTotal := 0
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
		UpdatedAt time.Time `json:"updated_at"`
	}

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/pulls?state=closed&sort=updated&direction=desc&per_page=100&page=%d",
			owner,
			repo,
			page,
		)

		l.Info().Msg(fmt.Sprintf("GitHub request URL: %s", url))

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "Team-Pulse-Metrics-Engine")

		resp, err := client.Do(req)
		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, fmt.Errorf("GitHub API Error: %s\n%s", resp.Status, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		var pullrequest PullRequestResponse

		if err := json.Unmarshal(body, &pullrequest); err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		returnedTotal += len(pullrequest)

		if len(pullrequest) == 0 {
			break
		}

		reachedOldPRs := false
		for _, pr := range pullrequest {
			if !since.IsZero() && (pr.UpdatedAt.Before(since) || pr.UpdatedAt.Equal(since)) {
				reachedOldPRs = true
				break
			}

			if pr.MergedAt == nil {
				skipped++
				continue
			}

			creator, err := queries.GetUserByGithubUsername(pr.User.Login)
			if err != nil {
				skipped++
				continue
			}

			actionUser := creator

			if pr.MergedBy != nil {
				if merger, err := queries.GetUserByGithubUsername(pr.MergedBy.Login); err == nil {
					actionUser = merger
				}
			}

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
				continue
			}

			imported++
		}

		if len(pullrequest) < 100 || reachedOldPRs {
			break
		}

		page++
	}

	l.Info().Msg(fmt.Sprintf("Pull requests returned by GitHub: %d", returnedTotal))
	l.Info().Msg(fmt.Sprintf("New pull requests inserted: %d", imported))

	return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, nil
}

// --------------------------------------------
// Issue History Sync
// --------------------------------------------
func SyncIssue(cfg GithubSyncConfig, since time.Time) (SyncCountSummary, error) {
	l := middleware.LogGet()
	owner := cfg.Owner
	repo := cfg.Repo
	token := cfg.Token
	client := cfg.Client

	page := 1
	returnedTotal := 0
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
		UpdatedAt time.Time  `json:"updated_at"`
	}

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/issues?state=all&sort=updated&direction=desc&per_page=100&page=%d",
			owner,
			repo,
			page,
		)
		if !since.IsZero() {
			url += fmt.Sprintf("&since=%s", since.UTC().Format(time.RFC3339))
		}

		l.Info().Msg(fmt.Sprintf("GitHub request URL: %s", url))

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "Team-Pulse-Metrics-Engine")

		resp, err := client.Do(req)
		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, fmt.Errorf("GitHub API Error: %s\n%s", resp.Status, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		var issues IssueResponse

		if err := json.Unmarshal(body, &issues); err != nil {
			return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, err
		}

		returnedTotal += len(issues)

		if len(issues) == 0 {
			break
		}

		reachedOldIssues := false
		for _, issue := range issues {
			if !since.IsZero() && (issue.UpdatedAt.Before(since) || issue.UpdatedAt.Equal(since)) {
				reachedOldIssues = true
				break
			}

			if issue.PullRequest != nil {
				continue
			}

			creator, err := queries.GetUserByGithubUsername(issue.User.Login)
			if err != nil {
				skipped++
				continue
			}

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
				skipped++
				continue
			}

			imported++
		}

		if len(issues) < 100 || reachedOldIssues {
			break
		}

		page++
	}

	l.Info().Msg(fmt.Sprintf("Issues returned by GitHub: %d", returnedTotal))
	l.Info().Msg(fmt.Sprintf("New issues inserted: %d", imported))

	return SyncCountSummary{Returned: returnedTotal, Imported: imported, Skipped: skipped}, nil
}

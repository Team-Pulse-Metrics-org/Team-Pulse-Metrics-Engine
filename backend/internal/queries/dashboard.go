package queries

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/google/uuid"
)

type CommitTrendItem struct {
	Week    string `json:"week"`
	Commits int    `json:"commits"`
}

type TopContributorItem struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Commits int    `json:"commits"`
}

type RecentActivityItem struct {
	Timestamp  string          `json:"timestamp"`
	Developer  string          `json:"developer"`
	Type       string          `json:"type"`
	Repository string          `json:"repository"`
	Message    string          `json:"message"`
	Payload    json.RawMessage `json:"payload"`
}

func GetDashboardStats() (map[string]int, error) {
	counts := map[string]int{
		"git_commit":          0,
		"pull_request_closed": 0,
		"task_completed":      0,
		"open_issue":          0,
	}

	query := `
		SELECT type, COUNT(*) 
		FROM activities 
		GROUP BY type
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var activityType string
		var count int
		if err := rows.Scan(&activityType, &count); err != nil {
			return nil, err
		}
		counts[activityType] = count
	}

	return counts, nil
}

func GetCommitTrend() ([]CommitTrendItem, error) {
	query := `
		SELECT TO_CHAR(DATE_TRUNC('week', logged_at), 'YYYY-MM-DD') AS week, COUNT(*) AS commits
		FROM activities
		WHERE type = 'git_commit'
		GROUP BY week
		ORDER BY week ASC
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trend []CommitTrendItem
	for rows.Next() {
		var item CommitTrendItem
		if err := rows.Scan(&item.Week, &item.Commits); err != nil {
			return nil, err
		}
		trend = append(trend, item)
	}

	if trend == nil {
		trend = []CommitTrendItem{}
	}

	return trend, nil
}

func GetTopContributors() ([]TopContributorItem, error) {
	query := `
		SELECT a.user_id::text, COALESCE(u.first_name || ' ' || u.last_name, u.github_username, 'Unknown') AS name, COUNT(*) AS commits
		FROM activities a
		LEFT JOIN users u ON a.user_id = u.id
		WHERE a.type = 'git_commit'
		GROUP BY a.user_id, u.first_name, u.last_name, u.github_username
		ORDER BY commits DESC
		LIMIT 5
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contributors []TopContributorItem
	for rows.Next() {
		var item TopContributorItem
		if err := rows.Scan(&item.UserID, &item.Name, &item.Commits); err != nil {
			return nil, err
		}
		contributors = append(contributors, item)
	}

	if contributors == nil {
		contributors = []TopContributorItem{}
	}

	return contributors, nil
}

func GetRecentActivity() ([]RecentActivityItem, error) {
	query := `
		SELECT COALESCE(u.first_name || ' ' || u.last_name, u.github_username, 'Unknown') AS developer, a.type, a.payload, a.logged_at
		FROM activities a
		LEFT JOIN users u ON a.user_id = u.id
		ORDER BY a.logged_at DESC
		LIMIT 5
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var activities []RecentActivityItem
	for rows.Next() {
		var developer string
		var activityType string
		var payloadBytes []byte
		var loggedAt time.Time

		if err := rows.Scan(&developer, &activityType, &payloadBytes, &loggedAt); err != nil {
			return nil, err
		}

		// Parse payload to extract repository and message
		var payloadMap map[string]interface{}
		if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
			log.Println("Error unmarshaling activity payload:", err)
		}

		repo := "Unknown"
		if payloadMap != nil {
			if r, ok := payloadMap["repository"]; ok {
				if rStr, ok := r.(string); ok {
					repo = rStr
				} else if rMap, ok := r.(map[string]interface{}); ok {
					if name, ok := rMap["name"].(string); ok {
						repo = name
					}
				}
			}
		}

		message := "No message"
		if payloadMap != nil {
			// Commits
			if commits, ok := payloadMap["commits"].([]interface{}); ok && len(commits) > 0 {
				if firstCommit, ok := commits[0].(map[string]interface{}); ok {
					if msg, ok := firstCommit["message"].(string); ok {
						message = msg
					}
				}
			} else if msg, ok := payloadMap["message"].(string); ok {
				message = msg
			}

			// PRs or Issues titles
			if message == "No message" {
				if title, ok := payloadMap["title"].(string); ok {
					message = title
				} else if pr, ok := payloadMap["pull_request"].(map[string]interface{}); ok {
					if title, ok := pr["title"].(string); ok {
						message = title
					}
				} else if issue, ok := payloadMap["issue"].(map[string]interface{}); ok {
					if title, ok := issue["title"].(string); ok {
						message = title
					}
				}
			}
		}

		activities = append(activities, RecentActivityItem{
			Timestamp:  loggedAt.Format(time.RFC3339),
			Developer:  developer,
			Type:       activityType,
			Repository: repo,
			Message:    message,
			Payload:    json.RawMessage(payloadBytes),
		})
	}

	if activities == nil {
		activities = []RecentActivityItem{}
	}

	return activities, nil
}

func SyncGithubIssues(owner, repo, token string) error {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues?state=all&per_page=100", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github api error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type GHUser struct {
		Login string `json:"login"`
	}
	type GHIssue struct {
		Number      int        `json:"number"`
		Title       string     `json:"title"`
		State       string     `json:"state"`
		User        GHUser     `json:"user"`
		PullRequest interface{} `json:"pull_request"`
		CreatedAt   time.Time  `json:"created_at"`
		ClosedAt    *time.Time `json:"closed_at"`
		HTMLURL     string     `json:"html_url"`
	}

	var issues []GHIssue
	if err := json.Unmarshal(body, &issues); err != nil {
		return err
	}

	for _, issue := range issues {
		if issue.PullRequest != nil {
			continue // skip PRs
		}

		creator, err := GetUserByGithubUsername(issue.User.Login)
		var creatorID *string
		var creatorUsername string = issue.User.Login
		var userID uuid.UUID
		if err == nil && creator != nil {
			idStr := creator.ID.String()
			creatorID = &idStr
			creatorUsername = creator.GithubUsername
			userID = creator.ID
		} else {
			users, err := GetAllUsers()
			if err == nil && len(users) > 0 {
				userID = users[0].ID
			} else {
				continue // skip if we have no users
			}
		}

		var activityType models.ActivityType
		var loggedAt time.Time
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

		activityPayload := map[string]any{
			"repository":         repo,
			"issue_number":       issue.Number,
			"title":              issue.Title,
			"state":              issue.State,
			"created_by_user_id": creatorID,
			"action_by_user_id":  userID.String(),
			"created_by":         creatorUsername,
			"action_by":          creatorUsername,
			"developer":          creatorUsername,
			"created_at":         issue.CreatedAt,
			"closed_at":          issue.ClosedAt,
			"url":                issue.HTMLURL,
		}

		payloadJSON, err := json.Marshal(activityPayload)
		if err != nil {
			continue
		}

		existing, err := FindIssueActivity(issue.Number, repo, repo)
		if err == nil && existing != nil {
			if existing.Type != activityType {
				existing.Type = activityType
				existing.Payload = payloadJSON
				existing.LoggedAt = loggedAt
				_ = UpdateActivity(*existing)
			}
		} else {
			newActivity := models.Activities{
				UserID:   userID,
				Type:     activityType,
				Payload:  payloadJSON,
				LoggedAt: loggedAt,
			}
			_ = CreateActivity(newActivity)
		}
	}

	return nil
}


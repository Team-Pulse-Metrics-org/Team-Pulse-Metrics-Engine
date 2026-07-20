package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type LastSync struct {
    LastSynced string `json:"last_synced"`
}

type GistResponse struct {
	Files map[string]struct {
		Content string `json:"content"`
	} `json:"files"`
}

type GistFile struct {
	Content string `json:"content"`
}

type GistRequest struct {
	Files map[string]GistFile `json:"files"`
}

func getGithubToken() string {
	token := os.Getenv("GITHUB_PAT")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	return token
}

func UpdateLastSyncGist(t time.Time) error {
	Sync := LastSync{
		LastSynced: t.Format(time.RFC3339),
	}

	SyncData, err := json.Marshal(Sync)
	if err != nil {
		return err
	}

	GistReq := GistRequest{
		Files: map[string]GistFile{
			"last_sync.json": {
				Content: string(SyncData),
			},
		},
	}

	body, err := json.Marshal(GistReq)
	if err != nil {
		return err
	}

	gistID := os.Getenv("GIST_ID")
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)

	tokens := []string{os.Getenv("GITHUB_TOKEN"), os.Getenv("GITHUB_PAT")}
	var lastErr error

	for _, token := range tokens {
		if token == "" {
			continue
		}

		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Team-Pulse-Metrics-Engine")

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		lastErr = fmt.Errorf("github gist update failed (%d): %s", resp.StatusCode, string(respBody))
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("no valid GitHub token found for Gist update")
}

func ReadLastSyncGist() (LastSync, error) {
	gistID := os.Getenv("GIST_ID")
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)

	tokens := []string{os.Getenv("GITHUB_TOKEN"), os.Getenv("GITHUB_PAT")}
	var lastErr error

	for _, token := range tokens {
		if token == "" {
			continue
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return LastSync{}, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "Team-Pulse-Metrics-Engine")

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("github gist read failed (%d): %s", resp.StatusCode, string(body))
			continue
		}

		var gistResp GistResponse
		if err := json.Unmarshal(body, &gistResp); err != nil {
			return LastSync{}, err
		}

		content, ok := gistResp.Files["last_sync.json"]
		if !ok {
			return LastSync{}, fmt.Errorf("last_sync.json file not found in gist")
		}

		var sync LastSync
		if err := json.Unmarshal([]byte(content.Content), &sync); err != nil {
			return LastSync{}, err
		}

		return sync, nil
	}

	if lastErr != nil {
		return LastSync{}, lastErr
	}
	return LastSync{}, fmt.Errorf("no valid GitHub token found for Gist read")
}
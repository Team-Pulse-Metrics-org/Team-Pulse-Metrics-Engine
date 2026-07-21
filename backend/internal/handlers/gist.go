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

func (h *MetricsHandler) getGithubToken() string {
	token := h.cfg.GithubPAT
	if token == "" {
		token = h.cfg.GithubToken
	}
	return token
}

func (h *MetricsHandler) UpdateLastSyncGist(t time.Time) error {
	Sync := LastSync{
		LastSynced: t.Format(time.RFC3339),
	}

	SyncData, err := json.Marshal(Sync)
	if err != nil {
		h.log.Error().Err(err).Msg("failed to marshal last sync timestamp")
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
		h.log.Error().Err(err).Msg("failed to marshal gist request payload")
		return err
	}

	gistID := h.cfg.GistID
	if gistID == "" {
		err := fmt.Errorf("GIST_ID is not configured")
		h.log.Error().Err(err).Msg("missing Gist configuration")
		return err
	}
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)

	tokens := []string{h.cfg.GithubToken, h.cfg.GithubPAT}
	var lastErr error

	for _, token := range tokens {
		if token == "" {
			continue
		}

		req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
		if err != nil {
			h.log.Error().Err(err).Msg("failed to create http request for gist update")
			return err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Team-Pulse-Metrics-Engine")

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			h.log.Warn().Err(err).Msg("failed request to GitHub Gist API")
			lastErr = err
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			h.log.Info().Time("last_synced", t).Msg("successfully updated last sync timestamp.")
			return nil
		}

		lastErr = fmt.Errorf("github gist update failed (%d): %s", resp.StatusCode, string(respBody))
		h.log.Warn().
			Int("status_code", resp.StatusCode).
			Str("body", string(respBody)).
			Msg("GitHub Gist update rejected")
	}

	if lastErr != nil {
		h.log.Error().Err(lastErr).Msg("all GitHub token attempts failed for Gist update")
		return lastErr
	}
	err = fmt.Errorf("no valid GitHub token found for Gist update")
	h.log.Error().Err(err).Msg("failed Gist update execution")
	return err
}

func (h *MetricsHandler) ReadLastSyncGist() (LastSync, error) {
	gistID := h.cfg.GistID
	if gistID == "" {
		err := fmt.Errorf("GIST_ID is not configured")
		h.log.Error().Err(err).Msg("missing Gist configuration")
		return LastSync{}, err
	}

	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)
	tokens := []string{os.Getenv("GITHUB_TOKEN"), os.Getenv("GITHUB_PAT")}
	var lastErr error

	for _, token := range tokens {
		if token == "" {
			continue
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			h.log.Error().Err(err).Msg("failed to create http request for gist read")
			return LastSync{}, err
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("User-Agent", "Team-Pulse-Metrics-Engine")

		client := http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			h.log.Warn().Err(err).Msg("failed request to read GitHub Gist, trying next token if available")
			lastErr = err
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("github gist read failed (%d): %s", resp.StatusCode, string(body))
			h.log.Warn().
				Int("status_code", resp.StatusCode).
				Str("body", string(body)).
				Msg("GitHub Gist read rejected")
			continue
		}

		var gistResp GistResponse
		if err := json.Unmarshal(body, &gistResp); err != nil {
			h.log.Error().Err(err).Msg("failed to unmarshal Gist response body")
			return LastSync{}, err
		}

		content, ok := gistResp.Files["last_sync.json"]
		if !ok {
			err := fmt.Errorf("last_sync.json file not found in gist")
			h.log.Error().Err(err).Msg("missing file in Gist payload")
			return LastSync{}, err
		}

		var sync LastSync
		if err := json.Unmarshal([]byte(content.Content), &sync); err != nil {
			h.log.Error().Err(err).Msg("failed to parse last_sync.json content")
			return LastSync{}, err
		}

		h.log.Debug().Str("last_synced", sync.LastSynced).Msg("successfully fetched last sync info from Gist")
		return sync, nil
	}

	if lastErr != nil {
		h.log.Error().Err(lastErr).Msg("all GitHub token attempts failed for Gist read")
		return LastSync{}, lastErr
	}
	err := fmt.Errorf("no valid GitHub token found for Gist read")
	h.log.Error().Err(err).Msg("failed Gist read execution")
	return LastSync{}, err
}

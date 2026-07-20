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

func UpdateLastSyncGist(t time.Time) error{
	Sync:=LastSync{
		LastSynced: t.Format(time.RFC3339),
	}

	SyncData, err := json.Marshal(Sync)
	if err!=nil{
		return err
	}

	GistReq:=GistRequest{
		Files: map[string]GistFile{
			"last_sync.json":{
				Content: string(SyncData),
			},
		},
	}

	body,err:=json.Marshal(GistReq)
	if err != nil {
		return err
	}

	gistID := os.Getenv("GIST_ID")
	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body),)
		if err != nil {
			panic(err)
		}

	token:=os.Getenv("GITHUB_TOKEN")

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")

	client:=http.Client{}
	resp,err:=client.Do(req)
	if err!=nil{
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("github gist update failed: %s", string(body))
	}

	return nil

}

func ReadLastSyncGist() (LastSync, error){
	gistID := os.Getenv("GIST_ID")

	url := fmt.Sprintf("https://api.github.com/gists/%s", gistID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return LastSync{}, err
	}
	token:=os.Getenv("GITHUB_TOKEN")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client:=http.Client{}
	resp,err:=client.Do(req)

	if err!=nil{
		return LastSync{},err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return LastSync{}, fmt.Errorf("github api error: %s", string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return LastSync{}, err
	}
	var gistResp GistResponse

	err = json.Unmarshal(body, &gistResp)
	if err != nil {
		return LastSync{}, err
	}

	content := gistResp.Files["last_sync.json"].Content

	var sync LastSync

	err = json.Unmarshal([]byte(content), &sync)
	if err != nil {
		return LastSync{}, err
	}

	return sync,nil
}
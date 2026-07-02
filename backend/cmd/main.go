package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func getCommits(c *gin.Context) {
	url := "https://api.github.com/repos/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/commits"
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal("GITHUB_TOKEN is not set in the environment")
	}

	req, err := http.NewRequestWithContext(c.Request.Context(), "GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create outbound request"})
		return
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-Github-Api-Version", "2022-11-28")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "Failed to reach Github API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusBadGateway, gin.H{
			"error":       "Github API returned an error",
			"github_code": resp.StatusCode,
		})
		return
	}

	var commits []models.GitHubCommit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Github reponse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"count":   len(commits),
		"commits": commits,
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	// url := "https://api.github.com/repos/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/commits"
	// token := os.Getenv("GITHUB_TOKEN")
	// if token == "" {
	// 	log.Fatal("GITHUB_TOKEN is not set in the environment")
	// }

	r := gin.Default()

	r.GET("/api/commits", getCommits)

	if err := r.Run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}

}

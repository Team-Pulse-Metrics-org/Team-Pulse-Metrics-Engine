package main

import (
	"fmt"
	"log"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/gin-contrib/cors"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	database.ConnectDB()
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/api/v1/webhook/github", handlers.HandleWebhook)
	r.POST("/api/v1/auth/login", handlers.HandleGithubLogin)
	r.GET("/api/v1/activities", handlers.GetActivities)
	r.GET("/api/v1/dashboard", handlers.GetDashboard)
	if err := r.Run(); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}

// removed markers

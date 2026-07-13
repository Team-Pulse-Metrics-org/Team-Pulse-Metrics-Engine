package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/gin-contrib/cors"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	l := middleware.LogGet()
	log.Println("app environment", os.Getenv("APP_ENV"))
	log.Println(os.Getenv("LOG_LEVEL"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.ConnectDB()

	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())
	r.Use(cors.Default())

	r.POST("/api/v1/webhook/github", handlers.HandleWebhook)
	r.POST("/api/v1/auth/login", handlers.HandleGithubLogin)
	r.GET("/api/v1/activities", handlers.GetActivities)
	r.GET("/api/v1/dashboard", handlers.GetDashboard)

	l.Info().
		Str("port", port).Msgf("Starting Team Pulse Metrics Server on port '%s'", port)

	if err := r.Run(":" + port); err != nil {
		l.Fatal().Err(err).Msg("server failed to start")
	}
}

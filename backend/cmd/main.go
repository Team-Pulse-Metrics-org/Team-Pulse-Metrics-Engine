package main

import (
	"fmt"
	"os"
	"time"

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	database.ConnectDB()

	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())
	r.Use(cors.New(config))

	r.POST("/api/v1/webhook/github", handlers.HandleWebhook)
	r.POST("/api/v1/auth/login", handlers.HandleGithubLogin)

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/activities", handlers.GetActivities)
		protected.GET("/dashboard", handlers.GetDashboard)
	}
	r.GET("/api/v1/admin/users", handlers.GetUsers)
	r.PUT(
		"/api/v1/admin/users/:id/role",
		handlers.UpdateUserRole,
	)
	r.DELETE(
		"/api/v1/admin/users/:id",
		handlers.DeleteUser,
	)
	r.POST(
		"/api/v1/admin/users",
		handlers.CreateUser,
	)
	l.Info().
		Str("port", port).Msgf("Starting Team Pulse Metrics Server on port '%s'", port)

	if err := r.Run(":" + port); err != nil {
		l.Fatal().Err(err).Msg("server failed to start")
	}
}

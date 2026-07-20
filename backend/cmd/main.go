package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/worker"
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

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if os.Getenv("ENABLE_METRICS_TOGGLE") == "true" {
		go worker.StartMetricsWorker(ctx)
	} else {
		l.Info().Msg("Background metrics worker is disabled")
	}

	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())
	r.Use(cors.New(config))

	r.POST("/api/v1/webhook/github", handlers.HandleWebhook)
	r.POST("/api/v1/auth/login", handlers.HandleGithubLogin)
	r.GET("/api/v1/admin/users", handlers.GetUsers)
	r.PUT("/api/v1/admin/users/:id/role", handlers.UpdateUserRole)
	r.DELETE("/api/v1/admin/users/:id", handlers.DeleteUser)
	r.POST("/api/v1/admin/users", handlers.CreateUser)

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/activities", handlers.GetActivities)
		protected.GET("/dashboard", handlers.GetDashboard)
		protected.GET("/profile", handlers.GetGitHubProfile)
		protected.POST("/sync", handlers.HandleSync)
		protected.GET("/metrics", handlers.HandleMetrics)
		protected.GET("/last-sync", handlers.GetLastSync)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	l.Info().
		Str("port", port).
		Msgf("Starting Team Pulse Metrics Server on port '%s'", port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	<-ctx.Done()
	l.Info().Msg("ShutDown signal received. Cleaning up resources...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		l.Error().Err(err).Msg("Server forced to shutdown")
	}

	l.Info().Msg("Team Pulse Metrics App shut down cleanly.")
}

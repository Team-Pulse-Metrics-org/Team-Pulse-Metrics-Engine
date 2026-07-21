package handlers

import (
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Container struct {
	Admin   *AdminHandler
	Auth    *AuthHandler
	Webhook *WebhookHandler
	Metrics *MetricsHandler
	Users   *UserHandler
}

func NewContainer(q *queries.Queries, cfg *config.Config, log zerolog.Logger) *Container {
	return &Container{
		Admin:   NewAdminHandler(q, cfg, log),
		Auth:    NewAuthHandler(q, cfg, log),
		Webhook: NewWebhookHandler(q, cfg, log),
		Metrics: NewMetricsHandler(q, cfg, log),
		Users:   NewUserHandler(q, cfg, log),
	}
}

func (c *Container) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	// Public Unprotected Routes
	api.POST("/webhook/github", c.Webhook.HandleWebhook)
	api.POST("/auth/login", c.Auth.HandleGithubLogin)

	// Admin Routes
	admin := api.Group("/admin")
	{
		admin.GET("/users", c.Admin.GetUsers)
		admin.POST("/users", c.Admin.CreateUser)
		admin.PUT("/users/:id/role", c.Admin.UpdateUserRole)
		admin.DELETE("/users/:id", c.Admin.DeleteUser)
	}

	// Protected Session Routes
	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		// User & Profile
		protected.GET("/profile", c.Users.GetGitHubProfile)
		protected.GET("/teams", c.Users.GetTeams)

		// Metrics, Dashboard, Sync & Teams
		protected.GET("/activities", c.Metrics.GetActivities)
		protected.GET("/dashboard", c.Metrics.GetDashboard)
		protected.POST("/sync", c.Metrics.HandleSync)
		protected.GET("/metrics", c.Metrics.HandleMetrics)
		protected.GET("/last-sync", c.Metrics.GetLastSync)
		protected.GET("/metrics/users", c.Metrics.HandleMetricDropDown)
		protected.GET("/metrics/user/:id", c.Metrics.HandleUserMetrics)
	}
}

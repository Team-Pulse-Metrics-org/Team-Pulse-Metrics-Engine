package handlers

import (
	"net/http"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type WebhookHandler struct {
	q   *queries.Queries
	cfg *config.Config
	log zerolog.Logger
}

func NewWebhookHandler(q *queries.Queries, cfg *config.Config, log zerolog.Logger) *WebhookHandler {
	return &WebhookHandler{
		q:   q,
		cfg: cfg,
		log: log,
	}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	h.log.Info().Msg("Webhook received from GitHub")

	eventType := c.GetHeader("X-GitHub-Event")
	h.log.Debug().Str("event_type", eventType).Msg("Processing Webhook event")
	switch eventType {
	case "pull_request":
		h.HandlePullRequest(c)
	case "push":
		h.HandlePush(c)
	case "issues":
		h.HandleIssueRequest(c)
	case "ping":
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
		return
	default:
		h.log.Warn().Str("event_type", eventType).Msg("Unhandled Github event type")
		c.JSON(http.StatusOK, gin.H{"message": "event ignored"})
		return
	}
}

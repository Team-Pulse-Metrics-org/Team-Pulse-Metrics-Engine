package worker

import (
	"context"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/handlers"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
)

func RunSync() {
	handlers.RunSync()
}

func StartSyncWorker(ctx context.Context) {
	l := middleware.LogGet()
	l.Info().Msg("Background Sync worker initialized")
	l.Info().Msg("Running initial GitHub sync on boot...")

	go RunSync()

	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.Info().Msg("Stopping background sync worker...")
			return
		case <-ticker.C:
			l.Info().Msg("Starting scheduled 12-hour sync...")
			go RunSync()
		}
	}
}

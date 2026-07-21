package worker

import (
	"context"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/handlers"
	"github.com/rs/zerolog"
)

type SyncWorker struct {
	metricsHandler *handlers.MetricsHandler
	log            zerolog.Logger
}

func NewSyncWorker(mh *handlers.MetricsHandler, log zerolog.Logger) *SyncWorker {
	return &SyncWorker{
		metricsHandler: mh,
		log:            log,
	}
}

func (w *SyncWorker) Start(ctx context.Context) {
	w.log.Info().Msg("Background Sync worker initialized")
	w.log.Info().Msg("Running initial GitHub sync on boot...")

	go w.metricsHandler.RunSync()

	ticker := time.NewTicker(12 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.log.Info().Msg("Stopping background sync worker...")
			return
		case <-ticker.C:
			w.log.Info().Msg("Starting scheduled 12-hour sync...")
			go w.metricsHandler.RunSync()
		}
	}
}

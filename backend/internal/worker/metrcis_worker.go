package worker

import (
	"context"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
)

func StartMetricsWorker(ctx context.Context) {
	l := middleware.LogGet()
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	l.Info().Msg("Background Metrics worker initialized")
	l.Info().Msg("running initial metrics sync on boot...")

	bootCtx, bootCancel := context.WithTimeout(ctx, 2*time.Minute)
	if err := queries.CreateMetric(bootCtx); err != nil {
		l.Error().Err(err).Msg("Initial boot metrics sync failed...")
	} else {
		l.Info().Msg("Initial boot metrics sync completed.")
	}
	bootCancel()

	for {
		select {
		case <-ctx.Done():
			l.Info().Msg("stopping metrics worker...")
			return
		case <-ticker.C:
			l.Info().Msg("Starting scheduled metrics sync...")
			syncCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
			err := queries.CreateMetric(syncCtx)
			cancel()

			if err != nil {
				l.Error().Msgf("Scheduled ")
			}
		}
	}
}

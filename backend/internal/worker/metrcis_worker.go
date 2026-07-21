package worker

import (
	"context"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/rs/zerolog"
)

func StartMetricsWorker(ctx context.Context, q *queries.Queries, cfg *config.Config, log zerolog.Logger) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	log.Info().Msg("Background Metrics worker initialized")
	log.Info().Msg("running initial metrics sync on boot...")

	bootCtx, bootCancel := context.WithTimeout(ctx, 2*time.Minute)
	if err := q.CreateMetric(bootCtx); err != nil {
		log.Error().Err(err).Msg("Initial boot metrics sync failed...")
	} else {
		log.Info().Msg("Initial boot metrics sync completed.")
	}
	bootCancel()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("stopping metrics worker...")
			return
		case <-ticker.C:
			log.Info().Msg("Starting scheduled metrics sync...")
			syncCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
			err := q.CreateMetric(syncCtx)
			cancel()

			if err != nil {
				log.Error().Msgf("Scheduled ")
			}
		}
	}
}

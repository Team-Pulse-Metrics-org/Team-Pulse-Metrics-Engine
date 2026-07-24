package worker

import (
	"context"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
)

func StartMetricsWorker(ctx context.Context, q *queries.Queries, cfg *config.Config, log zerolog.Logger) {
	log.Info().Msg("Background Metrics worker initialized")
	log.Info().Msg("running initial metrics sync on boot...")

	bootCtx, bootCancel := context.WithTimeout(ctx, 2*time.Minute)
	if err := q.CreateMetric(bootCtx); err != nil {
		log.Error().Err(err).Msg("Initial boot metrics sync failed...")
	} else {
		log.Info().Msg("Initial boot metrics sync completed.")
	}
	bootCancel()

	c := cron.New(cron.WithLocation(time.Local))

	runScheduledSync := func() {
		log.Info().Msg("Starting scheduled metrics sync...")
		syncCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		if err := q.CreateMetric(syncCtx); err != nil {
			log.Error().Err(err).Msg("Scheduled metrics sync failed")
		} else {
			log.Info().Msg("Scheduled metrics sync completed successfully")
		}
	}

	cronSpec := "0 0 * * 0"

	_, err := c.AddFunc(cronSpec, runScheduledSync)
	if err != nil {
		log.Error().Err(err).Msg("Failed to schedule metrics cron job")
		return
	}

	c.Start()
	log.Info().Str("spec", cronSpec).Msg("Metrics cron worker scheduled successfully")

	<-ctx.Done()
	log.Info().Msg("Stopping metrics worker...")

	cronCtx := c.Stop()
	<-cronCtx.Done()
	log.Info().Msg("Metrics worker gracefully stopped")
}

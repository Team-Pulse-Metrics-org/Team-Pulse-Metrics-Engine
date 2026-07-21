package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/config"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/database"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/handlers"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/middleware"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/queries"
	"github.com/Sheikh-Fahad-Ahmed/Team-Pulse-Metrics-Engine/internal/worker"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	l := middleware.LogGet()

	db, err := database.ConnectDB(cfg.DatabaseURL, l)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	q := queries.New(db)

	container := handlers.NewContainer(q, cfg, l)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	syncWorker := worker.NewSyncWorker(container.Metrics, l)

	go syncWorker.Start(ctx)

	if cfg.EnableMetricWorker {
		go worker.StartMetricsWorker(ctx, q, cfg, l)
	} else {
		l.Info().Msg("Background metrics worker is disabled")
	}

	corsConfig := cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r := gin.New()
	r.Use(middleware.RequestLogger())
	r.Use(middleware.Recovery())
	r.Use(cors.New(corsConfig))

	container.RegisterRoutes(r)

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  2 * time.Minute,
	}

	l.Info().
		Str("port", cfg.Port).
		Str("env", cfg.AppEnv).
		Msgf("Starting Team Pulse Metrics Server on port '%s'", cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	<-ctx.Done()
	l.Info().Msg("Shutdown signal received. Cleaning up resources...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		l.Error().Err(err).Msg("Server forced to shutdown")
	}

	l.Info().Msg("Team Pulse Metrics App shut down cleanly.")
}

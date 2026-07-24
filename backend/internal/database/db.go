package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // 👈 The new driver registration!
	"github.com/rs/zerolog"
)

func ConnectDB(dbURL string, log zerolog.Logger) (*sql.DB, error) {
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL environment variable is empty")
	}

	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		log.Error().Err(err).Msg("failed to open pgx database connection")
		return nil, err
	}

	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	start := time.Now()
	if err = db.Ping(); err != nil {
		log.Error().Err(err).Msg("failed to ping PostgreSQL database")
		return nil, err
	}
	pingDuration := time.Since(start)

	log.Info().Dur("ping_ms", pingDuration).Msg("Connected to PostgreSQL using pgx driver successfully!")
	return db, nil
}

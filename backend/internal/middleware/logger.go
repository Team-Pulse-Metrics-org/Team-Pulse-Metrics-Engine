package middleware

import (
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	adapter "github.com/axiomhq/axiom-go/adapters/zerolog"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var once sync.Once
var log zerolog.Logger

func LogGet() zerolog.Logger {
	once.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel)
		}

		var output io.Writer

		if os.Getenv("APP_ENV") == "development" {
			output = zerolog.ConsoleWriter{
				Out:        os.Stdout,
				TimeFormat: time.RFC3339,
			}
		} else {
			// Production
			axiomWriter, err := adapter.New()
			if err != nil {
				// Fallback
				output = os.Stderr
			} else {
				output = io.MultiWriter(os.Stderr, axiomWriter)
			}
		}

		buildInfo, _ := debug.ReadBuildInfo()

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("go_version", buildInfo.GoVersion).
			Logger()
	})

	return log
}

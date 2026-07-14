package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestLogger() gin.HandlerFunc {
	l := LogGet()

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		reqID := c.GetHeader("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}
		c.Set("request_id", reqID)
		c.Writer.Header().Set("X-Request-ID", reqID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		event := l.Info()
		if status >= 500 {
			event = l.Error()
		} else if status >= 400 {
			event = l.Warn()
		}

		event.
			Str("request_id", reqID).
			Str("method", c.Request.Method).
			Str("path", path).
			Int("status", status).
			Str("ip", c.ClientIP()).
			Float64("latency_ms", float64(latency.Microseconds())/1000.0).
			Int("body_size", c.Writer.Size())

		if len(c.Errors) > 0 {
			event.Str("errors", c.Errors.Last().Err.Error())
		}

		if status > 400 {
			event.Str("user-agent", c.Request.UserAgent())
		}

		event.Msg("request handled")

	}
}

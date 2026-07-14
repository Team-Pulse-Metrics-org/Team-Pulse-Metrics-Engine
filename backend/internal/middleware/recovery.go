package middleware

import (
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	l := LogGet()

	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				l.Error().Interface("panic", err).Str("request_id", c.GetString("request_id")).Str("path", c.Request.URL.Path).Bytes("stack", debug.Stack()).Msg("panic recovered")
				c.AbortWithStatus(500)
			}
		}()
		c.Next()

	}
}

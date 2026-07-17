package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		c.Next()

		status := c.Writer.Status()
		duration := time.Since(start).Milliseconds()

		level := slog.LevelInfo
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.Log(c.Request.Context(), level, "request",
			"method", method,
			"path", path,
			"status", status,
			"duration_ms", duration,
			"ip", clientIP,
		)
	}
}

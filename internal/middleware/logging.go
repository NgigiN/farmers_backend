package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		contentLength := c.Request.ContentLength

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()

		log.Printf("[%s] %s %s | Status: %d | Duration: %v | IP: %s | Request Size: %d bytes",
			method,
			path,
			c.Request.Proto,
			statusCode,
			duration,
			clientIP,
			contentLength,
		)
	}
}

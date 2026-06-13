package middleware

import (
	"time"

	"osint-scraper/internal/logger"

	"github.com/gin-gonic/gin"
)

// Zerologger returns a Gin middleware that logs HTTP requests using the application's zerolog instance.
func Zerologger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		// Choose log severity based on status code
		event := log.Info()
		if status >= 500 {
			event = log.Error()
		} else if status >= 400 {
			event = log.Warn()
		}

		if len(c.Errors) > 0 {
			event.Interface("errors", c.Errors.Errors())
		}

		// Format path to include query params
		fullPath := path
		if query != "" {
			fullPath = path + "?" + query
		}

		event.
			Int("status", status).
			Str("method", method).
			Str("path", fullPath).
			Str("ip", clientIP).
			Dur("latency", latency).
			Msg("HTTP request")
	}
}

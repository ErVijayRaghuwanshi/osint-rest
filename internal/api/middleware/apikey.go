package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// AdminAPIKey returns a Gin middleware that validates the X-API-Key header
// against the ADMIN_API_KEY environment variable.
func AdminAPIKey() gin.HandlerFunc {
	return func(c *gin.Context) {
		expected := os.Getenv("ADMIN_API_KEY")
		if expected == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "admin API key not configured"})
			c.Abort()
			return
		}

		key := c.GetHeader("X-API-Key")
		if key == "" || key != expected {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or missing API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}

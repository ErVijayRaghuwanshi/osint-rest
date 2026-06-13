package middleware

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"osint-scraper/internal/cache"
	"osint-scraper/internal/header"

	"github.com/gin-gonic/gin"
)

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// CacheMiddleware returns a Gin middleware that caches GET responses based on platform configuration.
func CacheMiddleware(c cache.Cache, hm *header.Manager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Only cache GET requests
		if ctx.Request.Method != http.MethodGet {
			ctx.Next()
			return
		}

		// Skip Swagger documentation
		if strings.HasPrefix(ctx.Request.URL.Path, "/swagger") {
			ctx.Next()
			return
		}

		// Skip health check endpoint
		if ctx.Request.URL.Path == "/health" {
			ctx.Next()
			return
		}

		// Form the cache key
		key := "http:cache:" + ctx.Request.RequestURI

		// Check cache
		if data, ok := c.Get(key); ok {
			ctx.Header("X-Cache", "HIT")
			ctx.Data(http.StatusOK, "application/json; charset=utf-8", data)
			ctx.Abort()
			return
		}

		// Set default header to MISS for outbound responses
		ctx.Header("X-Cache", "MISS")

		// Wrap ResponseWriter to capture the response body
		bw := &bodyWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBuffer(nil),
		}
		ctx.Writer = bw

		ctx.Next()

		// Store response in cache on successful 200 OK responses
		if ctx.Writer.Status() == http.StatusOK {
			// Determine TTL dynamically based on the platform name in the URI path (e.g. /api/jaco/userinfo -> jaco)
			ttl := 5 * time.Minute // default fallback
			parts := strings.Split(strings.TrimPrefix(ctx.Request.URL.Path, "/"), "/")
			if len(parts) >= 2 && parts[0] == "api" {
				platform := parts[1]
				secs := hm.GetCacheTTL(platform)
				if secs > 0 {
					ttl = time.Duration(secs) * time.Second
				}
			}

			c.Set(key, bw.body.Bytes(), ttl)
		}
	}
}

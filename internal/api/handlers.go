package api

import (
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

// HealthCheck godoc
// @Summary      Health check
// @Description  Returns the health status of the API including runtime info
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Router       /health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	c.JSON(http.StatusOK, gin.H{
		"status":     "ok",
		"goroutines": runtime.NumGoroutine(),
		"alloc_mb":   mem.Alloc / 1024 / 1024,
		"go_version": runtime.Version(),
	})
}

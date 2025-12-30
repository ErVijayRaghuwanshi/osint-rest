package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

// HealthCheck godoc
//
//	@Summary		Check service health
//	@Description	Returns OK if service is running
//	@Tags			Health
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (h *Handlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}


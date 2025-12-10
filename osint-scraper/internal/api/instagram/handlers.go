package instagram

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type Handler struct {
    svc *Service
}

func NewHandler(svc *Service) *Handler {
    return &Handler{svc: svc}
}

// Ping godoc
// @Summary      Check Instagram availability
// @Description  Pings instagram.com to ensure site is reachable
// @Tags         Instagram
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      503  {object}  map[string]string
// @Router       /api/instagram/ping [get]
func (h *Handler) Ping(c *gin.Context) {
    ok := h.svc.CheckWebsite()

    if ok {
        c.JSON(http.StatusOK, gin.H{"message": "instagram pong"})
    } else {
        c.JSON(http.StatusServiceUnavailable, gin.H{"message": "instagram unreachable"})
    }
}

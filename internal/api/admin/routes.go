package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"osint-scraper/internal/header"
)

func RegisterRoutes(rg *gin.RouterGroup, hm *header.Manager, log zerolog.Logger, headerFilePath string) {
	h := NewHandler(hm, log, headerFilePath)

	rg.GET("", h.ListPlatforms)
	rg.GET("/:platform", h.GetPlatformHeaders)
	rg.POST("/:platform", h.AddHeaderSet)
	rg.PUT("/:platform/:id", h.UpdateHeaderSet)
	rg.DELETE("/:platform/:id", h.DeleteHeaderSet)
	rg.POST("/:platform/:id/toggle", h.ToggleHeaderSet)
}

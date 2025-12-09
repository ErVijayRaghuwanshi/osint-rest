package api

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "osint-scraper/internal/api/snapchat"
    "osint-scraper/internal/config"
    "osint-scraper/internal/logger"
)

func NewRouter(cfg config.Config, log logger.Logger) *gin.Engine {
    r := gin.Default()

    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    h := NewHandlers()
    r.GET("/health", h.HealthCheck)
    r.GET("/ping", h.Ping) // optional global ping

    // -------------------------------
    // Snapchat Route Group
    // -------------------------------
    snapService := snapchat.NewService()
    snapGroup := r.Group("/api/snapchat") // <--- full prefix
    {
        snapchat.RegisterRoutes(snapGroup, snapService)
    }

    return r
}

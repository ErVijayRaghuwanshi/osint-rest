package api

import (
    "github.com/gin-gonic/gin"
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    "osint-scraper/internal/api/snapchat"
    "osint-scraper/internal/api/instagram"
    "osint-scraper/internal/api/x"
    "osint-scraper/internal/config"
    "osint-scraper/internal/logger"
)

func NewRouter(cfg config.Config, log logger.Logger) *gin.Engine {
    r := gin.Default()

    // Swagger endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    h := NewHandlers()
    r.GET("/health", h.HealthCheck)

    // -------------------------------
    // Snapchat Route Group
    // -------------------------------
    snapService := snapchat.NewService()
    snapGroup := r.Group("/api/snapchat") // <--- full prefix
    {
        snapchat.RegisterRoutes(snapGroup, snapService)
    }

    // -------------------------------
    // Instagram Route Group
    // -------------------------------
    instaService := instagram.NewService()
    instaGroup := r.Group("/api/instagram") // <--- full prefix
    {
        instagram.RegisterRoutes(instaGroup, instaService)
    }

    // -------------------------------
    // X Route Group
    // -------------------------------
    xService := x.NewService()
    xGroup := r.Group("/api/x") // <--- full prefix
    {
        x.RegisterRoutes(xGroup, xService)
    }

    return r
}

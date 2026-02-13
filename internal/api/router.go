package api

import (
	"time"

	"osint-scraper/internal/api/admin"
	"osint-scraper/internal/api/instagram"
	"osint-scraper/internal/api/jaco"
	"osint-scraper/internal/api/middleware"
	"osint-scraper/internal/api/snapchat"
	"osint-scraper/internal/api/x"
	"osint-scraper/internal/cache"
	"osint-scraper/internal/config"
	"osint-scraper/internal/header"
	"osint-scraper/internal/logger"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(cfg config.Config, log logger.Logger, hm *header.Manager, appCache cache.Cache, ht *header.HealthTracker, headerFilePath string) *gin.Engine {
	r := gin.Default()

	// CORS
	r.Use(middleware.CORS())

	// Rate limiter: 60 req/sec per IP, burst 100
	limiter := middleware.NewRateLimiter(60, 100, time.Second)
	r.Use(limiter.Middleware())

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h := NewHandlers()
	r.GET("/health", h.HealthCheck)

	// -------------------------------
	// Snapchat Route Group
	// -------------------------------
	snapService := snapchat.NewService(hm, log)
	snapService.SetCache(appCache)
	snapService.SetHealthTracker(ht)
	snapGroup := r.Group("/api/snapchat")
	{
		snapchat.RegisterRoutes(snapGroup, snapService, log)
	}

	// -------------------------------
	// Instagram Route Group
	// -------------------------------
	instaService := instagram.NewService(hm, log)
	instaService.SetCache(appCache)
	instaService.SetHealthTracker(ht)
	instaGroup := r.Group("/api/instagram")
	{
		instagram.RegisterRoutes(instaGroup, instaService, log)
	}

	// -------------------------------
	// X Route Group
	// -------------------------------
	xService := x.NewService(hm, log)
	xService.SetCache(appCache)
	xService.SetHealthTracker(ht)
	xGroup := r.Group("/api/x")
	{
		x.RegisterRoutes(xGroup, xService, log)
	}

	// -------------------------------
	// Jaco Route Group
	// -------------------------------
	jacoService := jaco.NewService(hm, log)
	jacoService.SetCache(appCache)
	jacoService.SetHealthTracker(ht)
	jacoGroup := r.Group("/api/jaco")
	{
		jaco.RegisterRoutes(jacoGroup, jacoService, log)
	}

	// -------------------------------
	// Admin Route Group (API key protected)
	// -------------------------------
	adminGroup := r.Group("/admin/headers", middleware.AdminAPIKey())
	{
		admin.RegisterRoutes(adminGroup, hm, log, headerFilePath)
	}

	return r
}

package api

import (
	"time"

	"osint-scraper/internal/api/admin"
	"osint-scraper/internal/api/middleware"
	"osint-scraper/internal/cache"
	"osint-scraper/internal/config"
	"osint-scraper/internal/header"
	"osint-scraper/internal/logger"
	"osint-scraper/internal/platform"
	"osint-scraper/internal/platform/instagram"
	"osint-scraper/internal/platform/jaco"
	"osint-scraper/internal/platform/snapchat"
	"osint-scraper/internal/platform/telegram"
	"osint-scraper/internal/platform/x"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(cfg config.Config, log logger.Logger, hm *header.Manager, appCache cache.Cache, ht *header.HealthTracker, headerFilePath string) *gin.Engine {
	r := gin.New()

	// Logger middleware (using zerolog)
	r.Use(middleware.Zerologger(log))

	// Recovery middleware (to catch panics and return 500)
	r.Use(gin.Recovery())

	// CORS
	r.Use(middleware.CORS())

	// Layered caching middleware
	r.Use(middleware.CacheMiddleware(appCache, hm))

	// Rate limiter: 60 req/sec per IP, burst 100
	limiter := middleware.NewRateLimiter(60, 100, time.Second)
	r.Use(limiter.Middleware())

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	h := NewHandlers()
	r.GET("/health", h.HealthCheck)

	// Initialize Platform Registry
	reg := platform.NewRegistry(log)

	// -------------------------------
	// Register Platforms
	// -------------------------------

	snapService := snapchat.NewService(hm, log)
	snapService.SetCache(appCache)
	snapService.SetHealthTracker(ht)
	reg.Register(snapService)

	instaService := instagram.NewService(hm, log)
	instaService.SetCache(appCache)
	instaService.SetHealthTracker(ht)
	reg.Register(instaService)

	xService := x.NewService(hm, log)
	xService.SetCache(appCache)
	xService.SetHealthTracker(ht)
	reg.Register(xService)

	jacoService := jaco.NewService(hm, log)
	jacoService.SetCache(appCache)
	jacoService.SetHealthTracker(ht)
	reg.Register(jacoService)

	teleService := telegram.NewService(hm, log)
	teleService.SetCache(appCache)
	teleService.SetHealthTracker(ht)
	reg.Register(teleService)

	// Mount all registered platforms under /api/<name>
	reg.MountAll(r)

	// -------------------------------
	// Admin Route Group (API key protected)
	// -------------------------------
	adminGroup := r.Group("/admin/headers", middleware.AdminAPIKey())
	{
		admin.RegisterRoutes(adminGroup, hm, log, headerFilePath)
	}

	return r
}

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"osint-scraper/internal/api"
	"osint-scraper/internal/cache"
	"osint-scraper/internal/config"
	"osint-scraper/internal/header"
	"osint-scraper/internal/logger"

	"github.com/joho/godotenv"

	"osint-scraper/docs" // Swagger docs
)

const (
	DefaultHeaderPath = "./config/headers.json"
	HeaderPathEnv     = "HEADERS_CONFIG_PATH"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	log := logger.New()

	// 1️⃣ Resolve header path
	headerPath := os.Getenv(HeaderPathEnv)
	if headerPath == "" {
		headerPath = DefaultHeaderPath
		log.Info().
			Str("path", headerPath).
			Msg("Using default header config")
		// log.Debug("HEADERS_CONFIG_PATH not set, using default", "path", headerPath)
	} else {
		log.Info().
			Str("path", headerPath).
			Msg("Using header config from env")
	}

	// ✅ Load .env (local/dev only)
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using system env")
	}
	cfg := config.Load()

	headerCfg, err := header.LoadHeaderConfig(headerPath)
	if err != nil {
		log.Error().Msgf("failed to load header config: %v", err)
	}

	headerManager := header.NewManager(headerCfg)

	// ✅ Dynamically set Swagger host
	if swaggerHost := os.Getenv("SWAGGER_HOST"); swaggerHost != "" {
		docs.SwaggerInfo.Host = swaggerHost
	}

	// ✅ Dynamically set version
	version := cfg.Version

	docs.SwaggerInfo.Version = version

	log.Info().
		Str("version", version).
		Msg("Starting OSINT Scraper API")

	// ✅ Dynamically set title
	docs.SwaggerInfo.Title = cfg.Title
	// ✅ Dynamically set description
	docs.SwaggerInfo.Description = cfg.Description

	// ✅ Initialize cache
	appCache := cache.NewMemoryCache(2 * time.Minute)
	defer appCache.Close()

	// ✅ Watch header config for hot reload
	watcherStop := make(chan struct{})
	go header.WatchAndReload(headerPath, headerManager, log, watcherStop)

	healthTracker := header.NewHealthTracker(headerManager, 5, log)

	router := api.NewRouter(*cfg, log, headerManager, appCache, healthTracker, headerPath)

	// ✅ Graceful shutdown
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	go func() {
		log.Info().Msg("OSINT Scraper API running on port " + cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")
	close(watcherStop)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited cleanly")
}

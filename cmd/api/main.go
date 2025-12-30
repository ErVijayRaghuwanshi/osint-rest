package main

import (
	"osint-scraper/internal/api"
	"osint-scraper/internal/config"
	"osint-scraper/internal/logger"

	_ "osint-scraper/docs" // Swagger docs
)

//	@title			OSINT Scraper API
//	@version		1.0
//	@description	High-performance OSINT scraping microservice.

// @host		localhost:8080
// @BasePath	/
func main() {
	log := logger.New()
	cfg := config.Load()

	router := api.NewRouter(cfg, log)

	log.Info().Msg("OSINT Scraper API running on port " + cfg.ServerPort)
	router.Run(":" + cfg.ServerPort)
}

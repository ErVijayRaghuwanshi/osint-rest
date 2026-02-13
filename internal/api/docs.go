package api

import (
	docs "osint-scraper/docs"
	"osint-scraper/internal/config"
)

// InitSwagger sets the Swagger host and version dynamically
func InitSwagger(cfg *config.Config) {
	docs.SwaggerInfo.Host = cfg.SwaggerHost
	docs.SwaggerInfo.Version = cfg.Version
}

package config

import (
	// "crypto/des"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	ServerPort  string
	SwaggerHost string
	Version     string
	Title       string
	Description string
}



// Load reads env vars and VERSION file
func Load() *Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	swaggerHost := os.Getenv("SWAGGER_HOST")
	if swaggerHost == "" {
		swaggerHost = "localhost:" + port
	}

	version := readVersion("VERSION")
	fmt.Printf("Loaded version: %s\n", version)

	title := os.Getenv("API_TITLE")
	if title == "" {
		title = "OSINT Scraper API"
	}

	description := os.Getenv("API_DESCRIPTION")
	if description == "" {
		description = "High-performance OSINT scraping microservice."
	}

	return &Config{
		ServerPort:  port,
		SwaggerHost: swaggerHost,
		Version:     version,
		Title:       title,
		Description: description,
	}
}

// readVersionFile reads the application version from env var or VERSION file
func readVersion(versionFilePath string) string {

	// 1️⃣ ENV override (best for Docker/K8s/CI)
	if v := strings.TrimSpace(os.Getenv("APP_VERSION")); v != "" {
		return v
	}

	// 2️⃣ VERSION file (local / repo-based)
	data, err := os.ReadFile(versionFilePath)
	if err == nil {
		if v := strings.TrimSpace(string(data)); v != "" {
			return v
		}
		log.Printf("WARNING: VERSION file empty, falling back to default")
	} else {
		log.Printf("WARNING: Could not read VERSION file: %v", err)
	}

	// 3️⃣ Hard fallback
	return "1.0.0"
}
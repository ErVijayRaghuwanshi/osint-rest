package config

import "os"

type Config struct {
	ServerPort string
}

func Load() Config {
	return Config{
		ServerPort: getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

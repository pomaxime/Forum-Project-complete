package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	DBPath      string
	SessionTTL  int // en heures
	Environment string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBPath:      getEnv("DB_PATH", "./database.db"),
		SessionTTL:  getEnvInt("SESSION_TTL_HOURS", 24),
		Environment: getEnv("ENV", "development"),
	}
}

// IsProd retourne true si l'application tourne en production.
func (c *Config) IsProd() bool {
	return c.Environment == "production"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}

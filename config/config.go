package config

import (
	"log/slog"
	"os"
	"strings"
)

// Application's configuration
type App struct {
	DBPath     string
	MQTTBroker string
	MQTTPort   string
	LogLevel   slog.Level
}

// Load reads configuration from environment variables with sensible defaults
func Load() *App {
	cfg := &App{
		DBPath:     getEnv("APP_DB_PATH", "environment_data.db"),
		MQTTBroker: getEnv("APP_MQTT_BROKER", "localhost"),
		MQTTPort:   getEnv("APP_MQTT_PORT", "1883"),
		LogLevel:   slog.LevelInfo, // Default log level
	}

	levelStr := getEnv("APP_LOG_LEVEL", "INFO")
	switch strings.ToUpper(levelStr) {
	case "DEBUG":
		cfg.LogLevel = slog.LevelDebug
	case "WARN":
		cfg.LogLevel = slog.LevelWarn
	case "ERROR":
		cfg.LogLevel = slog.LevelError
	}

	return cfg
}

// helper to read an environment variable or return a default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server  ServerConfig
	AI      AIConfig
	Logging LoggingConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// AIConfig holds AI provider configuration
type AIConfig struct {
	ProviderType string  // "openai", "ollama", "local"
	Model        string  // Model name
	Temperature  float64 // Temperature for generation
	MaxTokens    int     // Maximum tokens to generate
	BaseURL      string  // Base URL for API calls (for local providers)
	APIKey       string  // API key if required
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		},
		AI: AIConfig{
			ProviderType: getEnv("AI_PROVIDER", "openai"),
			Model:        getEnv("AI_MODEL", ""),
			Temperature:  getFloatEnv("AI_TEMPERATURE", 0.1),
			MaxTokens:    getIntEnv("AI_MAX_TOKENS", 1000),
			BaseURL:      getEnv("AI_BASE_URL", ""),
			APIKey:       getEnv("OPENAI_API_KEY", ""),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getIntEnv gets integer environment variable with fallback
func getIntEnv(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return fallback
}

// getFloatEnv gets float environment variable with fallback
func getFloatEnv(key string, fallback float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return fallback
}

// getDurationEnv gets duration environment variable with fallback
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return fallback
}

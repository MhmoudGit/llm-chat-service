package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	GroqAPIKey     string
	AppModel       string
	MaxTokens      int
	APIKey         string
	RateLimitRPS   int
	RateLimitBurst int
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		GroqAPIKey:     getEnv("GROQ_API_KEY", ""),
		AppModel:       getEnv("MODEL", "llama-3.3-70b-versatile"),
		MaxTokens:      getEnvInt("MAX_TOKENS", 1024),
		APIKey:         getEnv("API_KEY", ""),
		RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 10),   // Default 10 RPS
		RateLimitBurst: getEnvInt("RATE_LIMIT_BURST", 20), // Default burst 20
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return fallback
}

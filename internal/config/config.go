package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	GroqAPIKey string
	AppModel   string
	MaxTokens  int
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		GroqAPIKey: getEnv("GROQ_API_KEY", ""),
		AppModel:   getEnv("MODEL", "llama-3.3-70b-versatile"),
		MaxTokens:  getEnvInt("MAX_TOKENS", 1024),
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

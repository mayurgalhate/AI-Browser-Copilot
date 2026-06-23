package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	GroqAPIKey  string
	GroqBaseURL string
	DBURL       string
	RedisURL    string
	Port        string
}

var AppConfig *Config

func LoadConfig() {
	_ = godotenv.Load() // Ignore error if .env doesn't exist (e.g. in prod)

	AppConfig = &Config{
		GroqAPIKey:  os.Getenv("GROQ_API_KEY"),
		GroqBaseURL: getEnvOrDefault("GROQ_BASE_URL", "https://api.groq.com/openai/v1"),
		DBURL:       getEnvOrDefault("DB_URL", "postgres://user:password@localhost:5432/copilot?sslmode=disable"),
		RedisURL:    getEnvOrDefault("REDIS_URL", "localhost:6379"),
		Port:        getEnvOrDefault("PORT", "8080"),
	}

	if AppConfig.GroqAPIKey == "" {
		slog.Warn("GROQ_API_KEY is not set. LLM features will fail.")
	}
}

func getEnvOrDefault(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

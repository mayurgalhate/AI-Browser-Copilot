package db

import (
	"context"
	"log/slog"
	"os"

	"browser-copilot-backend/internal/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis() {
	redisURL := config.AppConfig.RedisURL

	RDB = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})

	if _, err := RDB.Ping(context.Background()).Result(); err != nil {
		slog.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}

	slog.Info("📦 Redis Cache Connected successfully!")
}

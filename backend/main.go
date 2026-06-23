package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"browser-copilot-backend/internal/api"
	"browser-copilot-backend/internal/config"
	"browser-copilot-backend/internal/db"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load centralized configuration
	config.LoadConfig()

	// Initialize Data Layers
	db.InitPostgres()
	db.InitRedis()

	// Setup Gin Router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Register Routes
	r.GET("/health", api.HealthHandler)
	r.GET("/ws", api.WsHandler)
	r.POST("/chat", api.ChatHandler)
	r.DELETE("/history/:session_id", api.HistoryClearHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// Start Server
	go func() {
		slog.Info("Server starting on port 8080...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("ListenAndServe Error", "error", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exiting")
}

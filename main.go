package main

import (
	"os"
	"log"
	"log/slog"
	"syscall"
	"context"
	"time"
	"os/signal"
	"net/http"
	"github.com/seeques/subman/internal/api"
	"github.com/seeques/subman/internal/config"
	"github.com/seeques/subman/internal/storage"
)

// @title Subscription Service API
// @version 1.0
// @description REST API for managing user subscriptions

// @host localhost:8080
// @BasePath /api/v1

func main() {
	cfg := config.LoadConfig()

	slog.Info("starting server", "port", cfg.Port)

	pool, err := storage.CreatePool()
	if err != nil {
		log.Fatalf("pgxpool creation failed: %v", err)
	}
	defer pool.Close()

	storage := storage.NewPostgresStorage(pool)

	s := api.NewServer(storage, cfg)

	go func() {
		if err := s.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("received shutdown signal")

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	slog.Info("server stopped")
}

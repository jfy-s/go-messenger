package main

import (
	"fmt"
	"log/slog"
	"messenger-auth/internal/config"
	"messenger-auth/internal/server"
	"messenger-auth/internal/storage/postgres"
	"os"
)

func main() {
	cfg := config.Load("config/config.yaml")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	storage, err := postgres.NewStorage(cfg.DatabaseUrl, logger.With("component", "storage"))
	if err != nil {
		panic(err)
	}
	defer storage.Close()
	srv := server.NewServer(cfg, logger.With("component", "server"), storage)
	logger.Info("starting auth service", "address", fmt.Sprintf("%s:%v", cfg.Hostname, cfg.Port))
	if err := srv.ServeHTTP(); err != nil {
		panic(err)
	}
}

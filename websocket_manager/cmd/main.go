package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"websocket_manager/internal/config"
	"websocket_manager/internal/server"
	"websocket_manager/internal/session"
	"websocket_manager/internal/storage/postgres"
)

func main() {
	ctx := context.Background()
	cfg := config.Load("config/config.yaml")
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	storage, err := postgres.NewStorage(cfg.DatabaseUrl, logger.With("component", "storage"))
	if err != nil {
		panic("failed to init storage")
	}
	logger.Debug("storage connected")
	defer storage.Close()
	logger.Debug("starting websocket server")
	hub := server.NewHub(ctx, storage, logger.With("component", "hub"))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		session.ServeWs(hub, w, r)
	})
	addr := fmt.Sprintf("%s:%v", cfg.Hostname, cfg.Port)
	logger.Info("websocket server started", "address", addr)
	err = http.ListenAndServe(addr, nil)

	if err != nil {
		logger.Error("failed to start server", "error", err)
	}
}

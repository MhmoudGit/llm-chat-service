package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"chat-service/internal/api"
	"chat-service/internal/chat"
	"chat-service/internal/config"
	"chat-service/internal/llm"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading .env file", "error", err)
		os.Exit(1)
	}

	cfg := config.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	history := chat.NewHistoryManager()
	llmClient := llm.NewClient(cfg)
	chatService := chat.NewService(history, llmClient)
	apiHandler := api.NewHandler(chatService)
	router := api.NewRouter(apiHandler, cfg)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}

	slog.Info("Server exited")
}

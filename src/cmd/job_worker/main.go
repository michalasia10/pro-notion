package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"src/internal/config"
	"src/internal/pkg/taskqueue"

	"github.com/hibiken/asynq"
)

func main() {
	cfg := config.Load()

	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisURL(),
		Password: cfg.Redis.Password,
	}

	mux := asynq.NewServeMux()

	// TODO: Register task handlers here
	// mux.HandleFunc("task:type", handlerFunc)

	server := taskqueue.NewServer(redisOpt, cfg.Async.Concurrency, cfg.Async.Queues)

	log.Println("Starting job worker...")

	if err := server.Start(mux); err != nil {
		log.Fatalf("server error: %v", err)
	}

	// Create context that listens for the interrupt signal from the OS
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Wait for interrupt signal
	<-ctx.Done()

	log.Println("Shutting down job worker...")

	// Graceful shutdown
	server.Shutdown()

	log.Println("Job worker stopped.")
}

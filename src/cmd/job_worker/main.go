package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"src/internal/config"
	"src/internal/pkg/taskqueue"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	redisOpt := asynq.RedisClientOpt{
		Addr:     cfg.RedisURL(),
		Password: cfg.Redis.Password,
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL(),
		Password: cfg.Redis.Password,
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis not reachable: %v", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("error closing redis client: %v", err)
		}
	}()

	mux := asynq.NewServeMux()

	// Register task handlers (stub handler keeps server valid; replace with real handlers)
	mux.HandleFunc("tasks:no-op", func(ctx context.Context, t *asynq.Task) error {
		return nil
	})

	server := taskqueue.NewServer(redisOpt, cfg.Async.Concurrency, cfg.Async.Queues)

	log.Println("Starting job worker...")

	// Start server in background so we can handle signals
	go func() {
		if err := server.Start(mux); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
	defer server.Shutdown()

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

package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"src/internal/config"
	"src/internal/pkg/eventbus"

	"github.com/ThreeDotsLabs/watermill"
)

func main() {
	config.Load()
	logger := watermill.NewStdLogger(false, false)
	router, err := eventbus.NewRouter(logger)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	// TODO: Register event subscribers here
	// router.AddHandler("handler_name", "topic_name", publisher, "handler_name", handlerFunc)

	log.Println("Starting event worker...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := router.Run(ctx); err != nil {
		log.Fatalf("router error: %v", err)
	}

	log.Println("Event worker stopped.")
}

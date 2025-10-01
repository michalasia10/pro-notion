package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"src/internal/config"
	sharedEvents "src/internal/modules/shared/domain/events"
	webhookEvents "src/internal/modules/webhooks/infrastructure/events"
	"src/internal/pkg/eventbus"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
)

func main() {
	config.Load()
	logger := watermill.NewStdLogger(false, false)
	publisher, err := eventbus.NewPublisher(logger)
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}

	router, err := eventbus.NewRouter(logger)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	// Register event subscribers
	triageLogger := log.New(log.Writer(), "webhook_triage: ", log.LstdFlags|log.Lshortfile)
	webhookTriage := webhookEvents.NewWebhookTriage(publisher, triageLogger)
	router.AddConsumerHandler(
		"webhook_triage",
		sharedEvents.NotionWebhookReceivedTopic,
		publisher.(message.Subscriber),
		webhookTriage.Handler(),
	)

	log.Println("Starting event worker...")

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := router.Run(ctx); err != nil {
		log.Fatalf("router error: %v", err)
	}

	log.Println("Event worker stopped.")
}

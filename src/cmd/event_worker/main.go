package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"src/internal/config"
	"src/internal/database"
	"src/internal/pkg/eventbus"

	projectsPostgres "src/internal/modules/projects/infrastructure/postgres"
	shared "src/internal/modules/shared/domain"
	sharedEvents "src/internal/modules/shared/domain/events"
	tasksPostgres "src/internal/modules/tasks/infrastructure/postgres"
	webhookEvents "src/internal/modules/webhooks/infrastructure/events"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/redis/go-redis/v9"
)

func main() {
	config.Load()
	logger := watermill.NewStdLogger(false, false)
	cfg := config.Get()
	eventBusCfg := eventbus.Config{
		Transport:       eventbus.Transport(cfg.EventBus.Transport),
		RedisOptions:    redis.Options{Addr: cfg.RedisURL(), Password: cfg.Redis.Password},
		ConsumerGroup:   cfg.EventBus.ConsumerGroup,
		ConsumerTimeout: cfg.EventBus.GetTimeoutInterval(),
		ClaimInterval:   cfg.EventBus.GetClaimInterval(),
	}

	pubSub, err := eventbus.NewPubSub(logger, eventBusCfg)
	if err != nil {
		log.Fatalf("failed to create publisher: %v", err)
	}
	defer func() {
		if pubSub.Close != nil {
			_ = pubSub.Close()
		}
	}()

	router, err := eventbus.NewRouter(logger)
	if err != nil {
		log.Fatalf("failed to create router: %v", err)
	}

	db := database.GormDB()
	tasksRepo := tasksPostgres.NewRepository(db)
	projectsRepo := projectsPostgres.NewProjectRepository(db)
	idGenerator := shared.NewUUIDGenerator()

	triageLogger := log.New(log.Writer(), "webhook_triage: ", log.LstdFlags|log.Lshortfile)
	webhookTriage := webhookEvents.NewWebhookTriage(
		pubSub.Publisher,
		triageLogger,
		tasksRepo,
		projectsRepo,
		idGenerator,
		nil,
	)

	router.AddConsumerHandler(
		"webhook_triage",
		sharedEvents.NotionWebhookReceivedTopic,
		pubSub.Subscriber,
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

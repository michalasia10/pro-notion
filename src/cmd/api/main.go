package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"src/internal/config"
	"src/internal/database"
	"src/internal/pkg/eventbus"
	"src/internal/pkg/notion"
	"src/internal/pkg/transaction"
	"src/internal/server"
	projectsPostgres "src/internal/modules/projects/infrastructure/postgres"
	shared "src/internal/modules/shared/domain"
	usersPostgres "src/internal/modules/users/infrastructure/postgres"

	_ "src/migrations"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/redis/go-redis/v9"
	"github.com/pressly/goose/v3"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	// Notify the main goroutine that the shutdown is complete
	done <- true
}

func main() {
	// Load configuration once at startup
	config.Load()
	cfg := config.Get()

	// Initialize DB (migrations are optional via RUN_MIGRATIONS=true)
	dbService := database.New()
	if strings.ToLower(os.Getenv("RUN_MIGRATIONS")) == "true" {
		if err := goose.SetDialect("postgres"); err != nil {
			panic(fmt.Sprintf("goose dialect error: %v", err))
		}
		if err := goose.Up(database.SQLDB(), "migrations"); err != nil {
			panic(fmt.Sprintf("failed to run migrations: %v", err))
		}
		log.Println("migrations applied")
	} else {
		log.Println("skipping migrations; set RUN_MIGRATIONS=true to apply on startup")
	}

	// Build dependencies
	gormDB := database.GormDB()
	txMgr := transaction.NewManager(gormDB)
	projectRepo := projectsPostgres.NewProjectRepository(gormDB)
	userRepo := usersPostgres.NewUserRepository(gormDB)
	idGen := shared.NewUUIDGenerator()
	clock := shared.NewSystemClock()
	notionService := notion.NewService(notion.ServiceConfig{
		ClientID:     cfg.Notion.ClientID,
		ClientSecret: cfg.Notion.ClientSecret,
		RedirectURI:  cfg.Notion.RedirectURL,
		APIVersion:   cfg.Notion.APIVersion,
	})

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL(),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	logger := watermill.NewStdLogger(false, false)
	eventBusCfg := eventbus.Config{
		Transport:     eventbus.Transport(cfg.EventBus.Transport),
		RedisOptions:  redis.Options{Addr: cfg.RedisURL(), Password: cfg.Redis.Password},
		ConsumerGroup: cfg.EventBus.ConsumerGroup,
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

	server := server.NewServer(server.Dependencies{
		DB:          dbService,
		Redis:       redisClient,
		Publisher:   pubSub.Publisher,
		TxMgr:       txMgr,
		ProjectRepo: projectRepo,
		UserRepo:    userRepo,
		IDGen:       idGen,
		Clock:       clock,
		Notion:      notionService,
	})

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	// Wait for the graceful shutdown to complete
	<-done
	log.Println("Graceful shutdown complete.")

	// Cleanup resources
	if err := redisClient.Close(); err != nil {
		log.Printf("error closing redis: %v", err)
	}
	database.Close()
}

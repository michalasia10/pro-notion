package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"

	"src/internal/config"
	"src/internal/database"
	"src/internal/pkg/eventbus"
)

type Server struct {
	port        int
	db          database.Service
	redisClient *redis.Client
	publisher   message.Publisher
}

func NewServer() *http.Server {
	cfg := config.Get()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL(),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	// Initialize Watermill publisher
	logger := watermill.NewStdLogger(false, false)
	publisher, err := eventbus.NewPublisher(logger)
	if err != nil {
		log.Fatalf("Failed to create Watermill publisher: %v", err)
	}

	serverInstance := &Server{
		port:        cfg.Port,
		redisClient: redisClient,
		db:          database.New(),
		publisher:   publisher,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverInstance.port),
		Handler:      serverInstance.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

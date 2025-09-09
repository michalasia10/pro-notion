package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"src/internal/config"
	"src/internal/database"
)

type Server struct {
	port        int
	db          database.Service
	redisClient *redis.Client
}

func NewServer() *http.Server {
	cfg := config.Get()

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL(),
		Password: cfg.Redis.Password,
		DB:       0,
	})

	serverInstance := &Server{
		port:        cfg.Port,
		redisClient: redisClient,
		db:          database.New(),
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

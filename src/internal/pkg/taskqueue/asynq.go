package taskqueue

import (
	"github.com/hibiken/asynq"
)

// NewClient creates a new Asynq client
func NewClient(redisOpt asynq.RedisConnOpt) *asynq.Client {
	return asynq.NewClient(redisOpt)
}

// NewServer creates a new Asynq server
func NewServer(redisOpt asynq.RedisConnOpt, concurrency int, queues map[string]int) *asynq.Server {
	return asynq.NewServer(redisOpt, asynq.Config{
		Concurrency: concurrency,
		Queues:      queues,
	})
}

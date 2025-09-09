package cache

import (
	"context"
	"time"
)

type Cache[T any] interface {
	Get(ctx context.Context, key string) (T, bool, error)
	Set(ctx context.Context, key string, data T, ttl time.Duration) error
	Clear(ctx context.Context) error
	Stats() Stats
}

type Stats struct {
	Hits   int64
	Misses int64
	Size   int64
}

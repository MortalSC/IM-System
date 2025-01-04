package repo

import (
	"context"
	"time"
)

type Cache interface {
	Put(ctx context.Context, key, val string, expire time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

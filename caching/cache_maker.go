package cache

import (
	"context"
	"errors"
	"fmt"
	"portfolio-profile-rpc/util"
	"time"

	"go.uber.org/zap"
)

var (
	ErrCacheMiss = errors.New("cache miss")
)

const (
	CacheTypeInMemory util.CacheType = "in_memory"
	CacheTypeRedis    util.CacheType = "redis"
)

// Repository pattern
// Cachier is an interface that defines the methods for interacting with a cache.
// It provides methods for setting, getting, adding to sets, checking if data is in a set,
// setting a value if it doesn't already exist, and deleting a key from the cache.
type Cachier interface {
	Set(ctx context.Context, key string, data any, ttl time.Duration) error
	Get(ctx context.Context, key string) (any, error)

	// Adds one or more values ​​to a set
	AddToSet(ctx context.Context, key string, data ...any) error
	IsDataInSet(ctx context.Context, key string, data any) (bool, error)

	// Using `SETNX`: SET if Not Exists. A.K.A `SET if Not Exists.`
	SetNX(ctx context.Context, key string, data any, ttl time.Duration) (bool, error)

	Del(ctx context.Context, key string) error
}

// Factory pattern
// NewCachierClient creates a new Cachier implementation based on the provided CacheConfig.
// It returns the Cachier implementation and an error if the cache type is not supported.
func NewCachierClient(
	cacheConfig util.CacheConfig,
	logger *zap.Logger,
) (Cachier, error) {
	switch cacheConfig.Type {
	case CacheTypeInMemory:
		return NewInMemoryClient(logger), nil

	case CacheTypeRedis:
		return NewRedisClient(cacheConfig, logger), nil

	default:
		return nil, fmt.Errorf("unsupported cache type: %s", cacheConfig.Type)
	}
}

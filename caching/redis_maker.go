package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"portfolio-profile-rpc/util"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// redisClient is a struct that holds a Redis client and a logger.
// It is used to interact with a Redis cache.
type redisClient struct {
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewRedisClient creates a new redisClient instance with the provided cache configuration and logger.
// The redisClient struct holds a Redis client and a logger, and is used to interact with a Redis cache.
// The Redis client is configured with the provided host, port, username, and password.
func NewRedisClient(cacheConfig util.CacheConfig, logger *zap.Logger) Cachier {
	addrString := fmt.Sprintf("%s:%d", cacheConfig.Host, cacheConfig.Port)

	return &redisClient{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     addrString,
			Username: cacheConfig.Username,
			Password: cacheConfig.Password,
		}),
		logger: logger,
	}
}

func (r redisClient) Set(ctx context.Context, key string, data any, ttl time.Duration) error {
	logger := util.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data)).
		With(zap.Duration("ttl", ttl))

		// Marshal data to byte slice
	dataBytes, err := json.Marshal(data) // Replace with suitable marshaling method
	if err != nil {
		logger.Info("failed to marshal data")
		return status.Error(codes.Internal, "failed to marshal data")
	}

	if err := r.redisClient.Set(ctx, key, dataBytes, ttl).Err(); err != nil {
		// logger.With(zap.Error(err)).Error("failed to set data into cache")
		logger.Info("failed to set data into cache: " + err.Error())
		return status.Error(codes.Internal, "failed to set data into cache")
	}

	return nil
}

func (r redisClient) Get(ctx context.Context, key string) (any, error) {
	logger := util.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key))

	data, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrCacheMiss
		}

		logger.With(zap.Error(err)).Error("failed to get data from cache")
		return nil, status.Error(codes.Internal, "failed to get data from cache")
	}

	return data, nil
}

func (r redisClient) AddToSet(ctx context.Context, key string, data ...any) error {
	logger := util.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data))

	if err := r.redisClient.SAdd(ctx, key, data...).Err(); err != nil {
		logger.With(zap.Error(err)).Error("failed to set data into set inside cache")
		return status.Error(codes.Internal, "failed to set data into set inside cache")
	}

	return nil
}

func (r redisClient) IsDataInSet(ctx context.Context, key string, data any) (bool, error) {
	logger := util.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data))

	result, err := r.redisClient.SIsMember(ctx, key, data).Result()
	if err != nil {
		logger.With(zap.Error(err)).Error("failed to check if data is member of set inside cache")
		return false, status.Error(codes.Internal, "failed to check if data is member of set inside cache")
	}

	return result, nil
}

func (r redisClient) SetNX(ctx context.Context, key string, data any, ttl time.Duration) (bool, error) {
	logger := util.LoggerWithContext(ctx, r.logger).
		With(zap.String("key", key)).
		With(zap.Any("data", data)).
		With(zap.Duration("ttl", ttl))

	// Marshal data to byte slice
	dataBytes, err := json.Marshal(data) // Replace with suitable marshaling method
	if err != nil {
		logger.Info("failed to marshal data")
		return false, status.Error(codes.Internal, "failed to marshal data")
	}

	ok, err := r.redisClient.SetNX(ctx, key, dataBytes, ttl).Result()
	if err != nil {
		logger.Info("failed to set data into cache: " + err.Error())
		return false, status.Error(codes.Internal, "failed to set data into cache")
	}

	return ok, nil
}

func (r redisClient) Del(ctx context.Context, key string) error {
	// Delete the key
	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		if err == redis.Nil {
			r.logger.Info("Key does not exist")
		} else {
			r.logger.Info("Error deleting key:" + err.Error())
		}
		return err
	}

	return nil
}

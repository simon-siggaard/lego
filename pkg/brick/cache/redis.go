package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// RedisClient is a Redis client.
type RedisClient struct {
	client  *redis.Client
	context context.Context
}

// NewRedisClient returns a new Redis client.
func NewRedisClient() *RedisClient {
	return &RedisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "passw0rd", // no password set
			DB:       0,          // use default DB
		}),
		context: context.Background(),
	}
}

// GET the cached value of key.
func (r *RedisClient) Get(key string) ([]byte, error) {
	return r.client.Get(r.context, key).Bytes()
}

// SET the cached value of key.
func (r *RedisClient) Set(key string, value []byte) error {
	return r.client.Set(r.context, key, value, 0).Err()
}

// Close the Redis client.
func (r *RedisClient) Close() {
	err := r.client.Close()
	if err != nil {
		panic(err)
	}
}

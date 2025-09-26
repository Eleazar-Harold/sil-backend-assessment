package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient handles Redis cache operations
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient creates a new Redis client
func NewRedisClient(addr, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisClient{
		client: rdb,
	}
}

// Set stores a value in Redis with an expiration time
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return c.client.Set(ctx, key, jsonValue, expiration).Err()
}

// Get retrieves a value from Redis
func (c *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("key not found")
		}
		return fmt.Errorf("failed to get value: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from Redis
func (c *RedisClient) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists checks if a key exists in Redis
func (c *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check key existence: %w", err)
	}
	return result > 0, nil
}

// SetNX sets a key only if it doesn't exist (atomic operation)
func (c *RedisClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return false, fmt.Errorf("failed to marshal value: %w", err)
	}

	result, err := c.client.SetNX(ctx, key, jsonValue, expiration).Result()
	if err != nil {
		return false, fmt.Errorf("failed to set key: %w", err)
	}

	return result, nil
}

// Expire sets an expiration time for a key
func (c *RedisClient) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL returns the time to live for a key
func (c *RedisClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// Increment increments a numeric value
func (c *RedisClient) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, key).Result()
}

// IncrementBy increments a numeric value by a specific amount
func (c *RedisClient) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, key, value).Result()
}

// Decrement decrements a numeric value
func (c *RedisClient) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, key).Result()
}

// DecrementBy decrements a numeric value by a specific amount
func (c *RedisClient) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.DecrBy(ctx, key, value).Result()
}

// ListPush adds elements to the left of a list
func (c *RedisClient) ListPush(ctx context.Context, key string, values ...interface{}) error {
	jsonValues := make([]interface{}, len(values))
	for i, v := range values {
		jsonValue, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value at index %d: %w", i, err)
		}
		jsonValues[i] = jsonValue
	}

	return c.client.LPush(ctx, key, jsonValues...).Err()
}

// ListPop removes and returns elements from the right of a list
func (c *RedisClient) ListPop(ctx context.Context, key string, dest interface{}) error {
	val, err := c.client.RPop(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("list is empty")
		}
		return fmt.Errorf("failed to pop from list: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

// ListLength returns the length of a list
func (c *RedisClient) ListLength(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}

// HashSet sets field-value pairs in a hash
func (c *RedisClient) HashSet(ctx context.Context, key string, values map[string]interface{}) error {
	jsonValues := make(map[string]interface{})
	for field, value := range values {
		jsonValue, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value for field %s: %w", field, err)
		}
		jsonValues[field] = jsonValue
	}

	return c.client.HMSet(ctx, key, jsonValues).Err()
}

// HashGet retrieves a field value from a hash
func (c *RedisClient) HashGet(ctx context.Context, key, field string, dest interface{}) error {
	val, err := c.client.HGet(ctx, key, field).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("field not found")
		}
		return fmt.Errorf("failed to get hash field: %w", err)
	}

	return json.Unmarshal([]byte(val), dest)
}

// HashGetAll retrieves all field-value pairs from a hash
func (c *RedisClient) HashGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HashDelete removes fields from a hash
func (c *RedisClient) HashDelete(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, key, fields...).Err()
}

// Ping tests the connection to Redis
func (c *RedisClient) Ping(ctx context.Context) error {
	return c.client.Ping(ctx).Err()
}

// Close closes the Redis connection
func (c *RedisClient) Close() error {
	return c.client.Close()
}

// GetClient returns the underlying Redis client for advanced operations
func (c *RedisClient) GetClient() *redis.Client {
	return c.client
}

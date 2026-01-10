package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Get retrieves a value by key. Returns redis.Nil error if key doesn't exist.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

// GetBytes retrieves a value as bytes (for binary data).
func (c *Client) GetBytes(ctx context.Context, key string) ([]byte, error) {
	return c.rdb.Get(ctx, key).Bytes()
}

// Exists checks if a key exists.
func (c *Client) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rdb.Exists(ctx, key).Result()
	return n > 0, err
}

// TTL gets the remaining time-to-live of a key.
// Returns -1 if key has no expiration, -2 if key doesn't exist.
func (c *Client) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.rdb.TTL(ctx, key).Result()
}

// MGet retrieves multiple values by keys.
func (c *Client) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return c.rdb.MGet(ctx, keys...).Result()
}

// Keys returns all keys matching a pattern. Use with caution in production.
func (c *Client) Keys(ctx context.Context, pattern string) ([]string, error) {
	return c.rdb.Keys(ctx, pattern).Result()
}

// IsNil checks if error is redis.Nil (key not found).
func IsNil(err error) bool {
	return err == redis.Nil
}

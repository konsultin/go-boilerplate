package redis

import (
	"context"
	"time"
)

// Set stores a key-value pair with optional expiration.
// Pass 0 for expiration to keep the key indefinitely.
func (c *Client) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

// SetNX sets a key only if it doesn't exist (atomic). Returns true if set, false if key exists.
func (c *Client) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.rdb.SetNX(ctx, key, value, expiration).Result()
}

// SetEX sets a key with mandatory expiration.
func (c *Client) SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.rdb.SetEx(ctx, key, value, expiration).Err()
}

// Del deletes one or more keys. Returns the number of keys deleted.
func (c *Client) Del(ctx context.Context, keys ...string) (int64, error) {
	return c.rdb.Del(ctx, keys...).Result()
}

// Expire sets expiration on a key. Returns false if key doesn't exist.
func (c *Client) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.rdb.Expire(ctx, key, expiration).Result()
}

// Incr increments a key's integer value by 1. Returns new value.
func (c *Client) Incr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Incr(ctx, key).Result()
}

// IncrBy increments a key's integer value by n. Returns new value.
func (c *Client) IncrBy(ctx context.Context, key string, n int64) (int64, error) {
	return c.rdb.IncrBy(ctx, key, n).Result()
}

// Decr decrements a key's integer value by 1. Returns new value.
func (c *Client) Decr(ctx context.Context, key string) (int64, error) {
	return c.rdb.Decr(ctx, key).Result()
}

// DecrBy decrements a key's integer value by n. Returns new value.
func (c *Client) DecrBy(ctx context.Context, key string, n int64) (int64, error) {
	return c.rdb.DecrBy(ctx, key, n).Result()
}

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client with helper methods.
type Client struct {
	rdb *redis.Client
}

// Config holds Redis connection configuration.
type Config struct {
	Host     string
	Port     int
	Password string
	DB       int

	// Connection Pool (optional, sensible defaults applied)
	PoolSize     int           // Max connections (default: 10 * GOMAXPROCS)
	MinIdleConns int           // Min idle connections (default: 0)
	DialTimeout  time.Duration // Connection timeout (default: 5s)
	ReadTimeout  time.Duration // Read timeout (default: 3s)
	WriteTimeout time.Duration // Write timeout (default: 3s)
}

// New creates a new Redis client wrapper.
func New(cfg Config) (*Client, error) {
	// Apply defaults
	if cfg.DialTimeout == 0 {
		cfg.DialTimeout = 5 * time.Second
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 3 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 3 * time.Second
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), cfg.DialTimeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

// Close closes the Redis connection.
func (c *Client) Close() error {
	if c == nil || c.rdb == nil {
		return nil
	}
	return c.rdb.Close()
}

// Raw returns the underlying redis.Client for advanced usage.
func (c *Client) Raw() *redis.Client {
	return c.rdb
}

// Ping checks if Redis is reachable.
func (c *Client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

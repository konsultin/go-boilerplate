package storage

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Client wraps the MinIO client with helper methods.
type Client struct {
	mc     *minio.Client
	bucket string
}

// Config holds MinIO/S3 connection configuration.
type Config struct {
	Endpoint  string // e.g., "localhost:9000" or "s3.amazonaws.com"
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
	Region    string
}

// New creates a new MinIO client wrapper.
func New(cfg Config) (*Client, error) {
	mc, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client init failed: %w", err)
	}

	return &Client{mc: mc, bucket: cfg.Bucket}, nil
}

// EnsureBucket creates the bucket if it doesn't exist.
func (c *Client) EnsureBucket(ctx context.Context) error {
	exists, err := c.mc.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("bucket check failed: %w", err)
	}
	if !exists {
		if err := c.mc.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("bucket creation failed: %w", err)
		}
	}
	return nil
}

// Bucket returns the configured bucket name.
func (c *Client) Bucket() string {
	return c.bucket
}

// Raw returns the underlying minio.Client for advanced usage.
func (c *Client) Raw() *minio.Client {
	return c.mc
}

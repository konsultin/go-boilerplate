package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
)

// UploadOptions configures upload behavior.
type UploadOptions struct {
	ContentType string            // MIME type (auto-detected if empty)
	Metadata    map[string]string // Custom metadata
}

// Upload uploads a file to storage. Returns the object path.
func (c *Client) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, opts *UploadOptions) (string, error) {
	putOpts := minio.PutObjectOptions{}
	if opts != nil {
		putOpts.ContentType = opts.ContentType
		putOpts.UserMetadata = opts.Metadata
	}

	// Auto-detect content type from extension if not provided
	if putOpts.ContentType == "" {
		putOpts.ContentType = detectContentType(objectName)
	}

	_, err := c.mc.PutObject(ctx, c.bucket, objectName, reader, size, putOpts)
	if err != nil {
		return "", fmt.Errorf("upload failed: %w", err)
	}

	return objectName, nil
}

// UploadFile is a convenience wrapper for uploading with auto-generated path.
// Returns: bucket/prefix/timestamp_filename
func (c *Client) UploadFile(ctx context.Context, prefix string, filename string, reader io.Reader, size int64, opts *UploadOptions) (string, error) {
	objectName := fmt.Sprintf("%s/%d_%s", prefix, time.Now().UnixMilli(), filename)
	return c.Upload(ctx, objectName, reader, size, opts)
}

// Delete removes an object from storage.
func (c *Client) Delete(ctx context.Context, objectName string) error {
	return c.mc.RemoveObject(ctx, c.bucket, objectName, minio.RemoveObjectOptions{})
}

// detectContentType returns MIME type based on file extension.
func detectContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".pdf":
		return "application/pdf"
	case ".json":
		return "application/json"
	case ".txt":
		return "text/plain"
	case ".html":
		return "text/html"
	case ".css":
		return "text/css"
	case ".js":
		return "application/javascript"
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	default:
		return "application/octet-stream"
	}
}

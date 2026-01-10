package storage

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// Download retrieves an object from storage.
func (c *Client) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := c.mc.GetObject(ctx, c.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("download failed: %w", err)
	}
	return obj, nil
}

// Stat returns object info without downloading.
func (c *Client) Stat(ctx context.Context, objectName string) (*minio.ObjectInfo, error) {
	info, err := c.mc.StatObject(ctx, c.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}
	return &info, nil
}

// Exists checks if an object exists.
func (c *Client) Exists(ctx context.Context, objectName string) (bool, error) {
	_, err := c.Stat(ctx, objectName)
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetPresignedURL generates a temporary download URL.
func (c *Client) GetPresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := c.mc.PresignedGetObject(ctx, c.bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("presign failed: %w", err)
	}
	return presignedURL.String(), nil
}

// GetPresignedUploadURL generates a temporary upload URL.
func (c *Client) GetPresignedUploadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	presignedURL, err := c.mc.PresignedPutObject(ctx, c.bucket, objectName, expiry)
	if err != nil {
		return "", fmt.Errorf("presign upload failed: %w", err)
	}
	return presignedURL.String(), nil
}

// List returns all objects with the given prefix.
func (c *Client) List(ctx context.Context, prefix string) ([]minio.ObjectInfo, error) {
	var objects []minio.ObjectInfo
	for obj := range c.mc.ListObjects(ctx, c.bucket, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		objects = append(objects, obj)
	}
	return objects, nil
}

package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strategic-insight-analyst/backend/config"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

var GCSClient *storage.Client

func InitializeGCS() error {
	ctx := context.Background()
	var client *storage.Client
	var err error

	credsFile := config.AppConfig.GoogleApplicationCredentials
	if credsFile != "" {
		// Use credentials file if provided (for local development)
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(credsFile))
	} else {
		// Otherwise, use Application Default Credentials (for production on GCP)
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return fmt.Errorf("failed to create gcs client: %w", err)
	}
	GCSClient = client
	return nil
}

func UploadFile(file multipart.File, header *multipart.FileHeader) (string, error) {
	bucketName := config.AppConfig.GCSBucketName
	if bucketName == "" {
		return "", fmt.Errorf("GCS_BUCKET_NAME environment variable not set in config")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), header.Filename)
	wc := GCSClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	return objectName, nil
}

func DownloadFile(objectName string) ([]byte, error) {
	bucketName := config.AppConfig.GCSBucketName
	if bucketName == "" {
		return nil, fmt.Errorf("GCS_BUCKET_NAME environment variable not set in config")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := GCSClient.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("NewReader: %v", err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("ReadAll: %v", err)
	}

	return data, nil
}

func DeleteFile(objectName string) error {
	bucketName := config.AppConfig.GCSBucketName
	if bucketName == "" {
		return fmt.Errorf("GCS_BUCKET_NAME environment variable not set in config")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	obj := GCSClient.Bucket(bucketName).Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("Delete: %v", err)
	}

	return nil
}

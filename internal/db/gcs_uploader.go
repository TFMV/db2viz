package db

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/storage"
)

type GCSUploader struct {
	bucketName string
	client     *storage.Client
}

func NewGCSUploader(cfg struct {
	ProjectID       string
	BucketName      string
	CredentialsFile string
}) (*GCSUploader, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	return &GCSUploader{
		bucketName: cfg.BucketName,
		client:     client,
	}, nil
}

func (uploader *GCSUploader) Upload(data []map[string]interface{}) error {
	ctx := context.Background()
	bucket := uploader.client.Bucket(uploader.bucketName)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	objectName := fmt.Sprintf("data-%d.json", time.Now().Unix())
	writer := bucket.Object(objectName).NewWriter(ctx)
	writer.ContentType = "application/json"

	if _, err := bytes.NewReader(jsonData).WriteTo(writer); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	log.Printf("Data successfully uploaded to GCS: %s", objectName)
	return nil
}

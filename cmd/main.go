package main

import (
	"context"
	"log"

	"github.com/TFMV/db2viz/config"
	"github.com/TFMV/db2viz/internal/data"
	"github.com/TFMV/db2viz/internal/db"
	"github.com/TFMV/db2viz/internal/gcp"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Postgres
	conn, err := db.ConnectPostgres(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer conn.Close(ctx)

	// Load data from Postgres
	loader := data.NewLoader(conn)
	records, err := loader.LoadData(ctx)
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Transform data
	transformer := data.NewTransformer()
	transformedData := transformer.Transform(records)

	// Upload data to BigQuery
	bqUploader := gcp.NewBigQueryUploader(cfg.BigQuery)
	err = bqUploader.Upload(ctx, transformedData)
	if err != nil {
		log.Fatalf("Failed to upload data to BigQuery: %v", err)
	}

	log.Println("Data pipeline completed successfully!")
}

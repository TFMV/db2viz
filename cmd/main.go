package main

import (
	"context"
	"log"

	"github.com/TFMV/db2viz/config"
	"github.com/TFMV/db2viz/internal/db"
	"github.com/TFMV/db2viz/internal/gcp"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to Postgres
	conn, err := db.NewPostgresConnection(cfg.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer conn.Close(ctx)

	// Load data from Postgres
	records, err := db.LoadData(conn, cfg.Postgres.Table)
	if err != nil {
		log.Fatalf("Failed to load data: %v", err)
	}

	// Publish data to Pub/Sub
	pubSubPublisher, err := gcp.NewPubSubPublisher(cfg.PubSub)
	if err != nil {
		log.Fatalf("Failed to create Pub/Sub publisher: %v", err)
	}

	err = pubSubPublisher.Publish(ctx, records)
	if err != nil {
		log.Fatalf("Failed to publish data to Pub/Sub: %v", err)
	}

	log.Println("Data pipeline completed successfully!")
}

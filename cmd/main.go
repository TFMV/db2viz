package main

import (
	"context"
	"log"
	"sync"

	"github.com/TFMV/db2viz/config"
	"github.com/TFMV/db2viz/internal/data"
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
	defer conn.Pool.Close()

	// Create a wait group for concurrency control
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, cfg.Postgres.Concurrency)

	for _, tableConfig := range cfg.Postgres.Tables {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(tableConfig config.TableConfig) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// Load data from Postgres
			loader := data.NewLoader(conn, tableConfig.Name)
			records, err := loader.LoadData(ctx)
			if err != nil {
				log.Printf("Failed to load data for table %s: %v", tableConfig.Name, err)
				return
			}

			// Transform data
			transformer := data.NewTransformer()
			transformedData := transformer.Transform(records)

			// Initialize PubSub client for the specific topic
			pubSubClient, err := gcp.NewPubSubClient(ctx, cfg.PubSub, tableConfig.TopicID)
			if err != nil {
				log.Printf("Failed to create PubSub client for table %s: %v", tableConfig.Name, err)
				return
			}

			// Publish data to PubSub
			err = pubSubClient.Publish(ctx, transformedData, cfg.PubSub.Workers)
			if err != nil {
				log.Printf("Failed to publish data for table %s: %v", tableConfig.Name, err)
			} else {
				log.Printf("Successfully published data for table %s", tableConfig.Name)
			}
		}(tableConfig)
	}

	wg.Wait()
	log.Println("Data pipeline completed successfully!")
}

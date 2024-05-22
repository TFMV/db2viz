package gcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"cloud.google.com/go/pubsub"
	"github.com/TFMV/db2viz/config"
	"google.golang.org/api/option"
)

type PubSubClient struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPubSubClient(ctx context.Context, cfg config.PubSubConfig, topicID string) (*PubSubClient, error) {
	client, err := pubsub.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.Credentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create PubSub client: %v", err)
	}

	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to check if topic exists: %v", err)
	}
	if !exists {
		topic, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic: %v", err)
		}
	}

	return &PubSubClient{
		client: client,
		topic:  topic,
	}, nil
}

func (p *PubSubClient) Publish(ctx context.Context, data []map[string]interface{}, workers int) error {
	recordsCh := make(chan map[string]interface{}, len(data))
	resultsCh := make(chan error, len(data))

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go p.worker(ctx, recordsCh, resultsCh, &wg)
	}

	// Send data to workers
	for _, record := range data {
		recordsCh <- record
	}
	close(recordsCh)

	// Wait for all workers to finish
	wg.Wait()
	close(resultsCh)

	// Collect results
	for err := range resultsCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PubSubClient) worker(ctx context.Context, recordsCh <-chan map[string]interface{}, resultsCh chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	for record := range recordsCh {
		jsonData, err := json.Marshal(record)
		if err != nil {
			resultsCh <- fmt.Errorf("failed to marshal record: %v", err)
			continue
		}

		result := p.topic.Publish(ctx, &pubsub.Message{
			Data: jsonData,
		})

		_, err = result.Get(ctx)
		if err != nil {
			resultsCh <- fmt.Errorf("failed to publish message: %v", err)
			continue
		}

		resultsCh <- nil
	}
}

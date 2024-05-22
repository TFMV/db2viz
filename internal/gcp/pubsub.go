package gcp

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/TFMV/db2viz/config"
	"google.golang.org/api/option"
)

type PubSubPublisher struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPubSubPublisher(cfg config.PubSubConfig) (*PubSubPublisher, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, cfg.ProjectID, option.WithCredentialsFile(cfg.Credentials))
	if err != nil {
		return nil, err
	}

	// Check if the topic exists, and create it if it doesn't
	topic := client.Topic(cfg.TopicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		topic, err = client.CreateTopic(ctx, cfg.TopicID)
		if err != nil {
			return nil, fmt.Errorf("failed to create topic: %v", err)
		}
		fmt.Printf("Created topic: %s\n", cfg.TopicID)
	} else {
		fmt.Printf("Using existing topic: %s\n", cfg.TopicID)
	}

	return &PubSubPublisher{
		client: client,
		topic:  topic,
	}, nil
}

func (publisher *PubSubPublisher) Publish(ctx context.Context, data []map[string]interface{}) error {
	for _, record := range data {
		message, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %v", err)
		}

		result := publisher.topic.Publish(ctx, &pubsub.Message{
			Data: message,
		})

		// Block until the result is returned and log server-assigned message ID
		id, err := result.Get(ctx)
		if err != nil {
			return fmt.Errorf("failed to publish message: %v", err)
		}
		fmt.Printf("Published a message with a message ID: %s\n", id)
	}

	return nil
}

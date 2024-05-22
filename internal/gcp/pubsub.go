package gcp

import (
	"context"
	"encoding/json"
	"fmt"

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

func (p *PubSubClient) Publish(ctx context.Context, data []map[string]interface{}) error {
	for _, record := range data {
		jsonData, err := json.Marshal(record)
		if err != nil {
			return fmt.Errorf("failed to marshal record: %v", err)
		}

		result := p.topic.Publish(ctx, &pubsub.Message{
			Data: jsonData,
		})

		_, err = result.Get(ctx)
		if err != nil {
			return fmt.Errorf("failed to publish message: %v", err)
		}
	}
	return nil
}

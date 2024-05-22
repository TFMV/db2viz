package gcp

import (
	"context"
	"log"

	"cloud.google.com/go/pubsub"
)

type PubSubClient struct {
	projectID string
	client    *pubsub.Client
}

func NewPubSubClient(cfg struct {
	ProjectID       string
	CredentialsFile string
}) (*PubSubClient, error) {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, err
	}
	return &PubSubClient{
		projectID: cfg.ProjectID,
		client:    client,
	}, nil
}

func (p *PubSubClient) Publish(topicName string, message []byte) error {
	ctx := context.Background()
	topic := p.client.Topic(topicName)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: message,
	})

	id, err := result.Get(ctx)
	if err != nil {
		return err
	}

	log.Printf("Published message to topic %s; message ID: %s", topicName, id)
	return nil
}

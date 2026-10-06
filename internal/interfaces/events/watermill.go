package events

import (
	"context"
	"encoding/json"

	"fmt"

	"github.com/ThreeDotsLabs/watermill"

	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	appconfig "github.com/zekielmp/Bitly/internal/config"
	"github.com/zekielmp/Bitly/internal/providers"
)

type EventPublisher struct {
	publisher message.Publisher
	queueName string
}

func (ep *EventPublisher) Publish(eventType string, payload interface{}, metadata map[string]string) error {

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), data)

	/* Add metadata */
	msg.Metadata.Set("event-type", eventType)
	for k, v := range metadata {
		msg.Metadata.Set(k, v)
	}
	return ep.publisher.Publish(ep.queueName, msg)
}

func (ep *EventPublisher) Close() error {
	return ep.publisher.Close()
}

func NewEventPublisher(ctx context.Context, cfg appconfig.AwsConfig) (*EventPublisher, error) {
	logger := watermill.NewStdLogger(false, true)

	awsCfg, err := providers.CreateAWSConfig(ctx, cfg.S3Endpoint, cfg.Region)
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS config: %w", err)
	}

	/*Create watermill SQS publisher config */
	publisherConfig := sqs.PublisherConfig{
		AWSConfig: awsCfg,
		Marshaler: nil,
	}

	/* Create the publisher with custom config */
	publisher, err := sqs.NewPublisher(publisherConfig, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create SQS publisher: %w", err)
	}

	return &EventPublisher{
		publisher: publisher,
		queueName: cfg.EventQueueName,
	}, nil
}

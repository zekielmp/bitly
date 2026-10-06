package events

import (
	"context"
	"encoding/json"

	// "fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/aws/aws-sdk-go-v2/aws"

	// "github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
	// "github.com/aws/smithy-go/endpoints"
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


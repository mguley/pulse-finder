package event

import (
	"domain"
	"encoding/json"
	"fmt"

	"github.com/nats-io/nats.go"
)

// NatsEventDispatcher uses NATS as the message broker to publish events to subscribers.
type NatsEventDispatcher struct {
	nc *nats.Conn // NATS connection for publishing messages.
}

// NewNatsEventDispatcher initializes a new NatsEventDispatcher with a NATS connection.
func NewNatsEventDispatcher(url string) (*NatsEventDispatcher, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS server: %w", err)
	}
	return &NatsEventDispatcher{nc: nc}, nil
}

// Dispatch publishes the specified event to a NATS topic based on the event type.
func (d *NatsEventDispatcher) Dispatch(e domain.Event) error {
	topic := fmt.Sprintf("event.%s", e.EventType())
	payload, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err = d.nc.Publish(topic, payload); err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}
	fmt.Printf("Published event to topic %s: %s\n", topic, string(payload))
	return nil
}

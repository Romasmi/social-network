package publisher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Romasmi/social-network/internal/events"
	"github.com/Romasmi/social-network/internal/infra/kafka_client"
	kafkago "github.com/confluentinc/confluent-kafka-go/kafka"
)

type Publisher interface {
	Publish(ctx context.Context, event *events.Event) error
	PublishToTopic(ctx context.Context, event *events.Event, topic string) error
	Flush(timeout time.Duration) error
}

type PublisherImpl struct {
	kafka        *kafka_client.Connection
	defaultTopic string
}

type PublisherConfig struct {
	DefaultTopic string
}

func NewPublisher(kafkaConn *kafka_client.Connection, config PublisherConfig) Publisher {
	if config.DefaultTopic == "" {
		config.DefaultTopic = "social-network-events"
	}

	return &PublisherImpl{
		kafka:        kafkaConn,
		defaultTopic: config.DefaultTopic,
	}
}

func (p *PublisherImpl) Publish(ctx context.Context, event *events.Event) error {
	return p.PublishToTopic(ctx, event, p.defaultTopic)
}

func (p *PublisherImpl) PublishToTopic(ctx context.Context, event *events.Event, topic string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("context cancelled: %w", err)
	}

	eventData, err := event.ToJSON()
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	key := fmt.Sprintf("%s:%s", event.Type, event.ID)

	message := &kafkago.Message{
		TopicPartition: kafkago.TopicPartition{
			Topic:     &topic,
			Partition: kafkago.PartitionAny,
		},
		Key:   []byte(key),
		Value: eventData,
		Headers: []kafkago.Header{
			{Key: "event_type", Value: []byte(event.Type)},
			{Key: "event_id", Value: []byte(event.ID)},
			{Key: "source", Value: []byte(event.Source)},
			{Key: "version", Value: []byte(event.Version)},
		},
		Timestamp: event.Timestamp,
	}

	deliveryChan := make(chan kafkago.Event, 1)
	err = p.kafka.Producer.Produce(message, deliveryChan)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}

	select {
	case e := <-deliveryChan:
		m := e.(*kafkago.Message)
		if m.TopicPartition.Error != nil {
			return fmt.Errorf("delivery failed: %w", m.TopicPartition.Error)
		}
		log.Printf("Published event: ID=%s, Type=%s, Topic=%s, Partition=%d, Offset=%d",
			event.ID, event.Type, topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	case <-time.After(5 * time.Second):
		return fmt.Errorf("delivery timeout after 5 seconds")
	case <-ctx.Done():
		return fmt.Errorf("context cancelled during delivery: %w", ctx.Err())
	}

	return nil
}

func (p *PublisherImpl) Flush(timeout time.Duration) error {
	remaining := p.kafka.Producer.Flush(int(timeout.Milliseconds()))
	if remaining > 0 {
		return fmt.Errorf("failed to flush %d messages within timeout", remaining)
	}
	log.Println("Publisher flushed successfully")
	return nil
}

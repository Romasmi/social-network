package events

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Romasmi/social-network/internal/infra/kafka"
	kafkago "github.com/confluentinc/confluent-kafka-go/kafka"
)

type ConsumerConfig struct {
	Topics            []string
	PollTimeout       time.Duration
	MaxRetries        int
	AutoCommit        bool
	ProcessingTimeout time.Duration
}

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
	Pause(partitions []kafkago.TopicPartition) error
	Resume(partitions []kafkago.TopicPartition) error
	Seek(partition kafkago.TopicPartition, timeoutMs int) error
	CommitOffsets(offsets []kafkago.TopicPartition) error
}

type ConsumerImpl struct {
	kafka    *kafka.Connection
	registry EventRegistry
	config   ConsumerConfig
	mu       sync.RWMutex
	running  bool
}

func NewConsumer(kafkaConn *kafka.Connection, registry EventRegistry, config ConsumerConfig) Consumer {
	if config.PollTimeout == 0 {
		config.PollTimeout = 5 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.ProcessingTimeout == 0 {
		config.ProcessingTimeout = 30 * time.Second
	}

	return &ConsumerImpl{
		kafka:    kafkaConn,
		registry: registry,
		config:   config,
		running:  false,
	}
}

func (c *ConsumerImpl) Start(ctx context.Context) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("consumer already running")
	}
	c.running = true
	c.mu.Unlock()

	err := c.kafka.Consumer.SubscribeTopics(c.config.Topics, nil)
	if err != nil {
		c.running = false
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}
	c.consumeLoop(ctx)
	return nil
}

func (c *ConsumerImpl) consumeLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("ConsumerImpl shutting down")
			return
		default:
			msg, err := c.kafka.Consumer.ReadMessage(c.config.PollTimeout)
			if err != nil {
				var kafkaErr kafkago.Error
				if errors.As(err, &kafkaErr) {
					if kafkaErr.Code() == kafkago.ErrTimedOut {
						continue
					}
				}
				log.Printf("ConsumerImpl error: %v", err)
				continue
			}

			if err := c.processMessageWithRetry(ctx, msg); err != nil {
				log.Printf("Failed to process message after retries: %v", err)
			}

			if !c.config.AutoCommit {
				if _, err := c.kafka.Consumer.CommitMessage(msg); err != nil {
					log.Printf("Failed to commit offset: %v", err)
				}
			}
		}
	}
}

func (c *ConsumerImpl) processMessageWithRetry(ctx context.Context, msg *kafkago.Message) error {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := 1 * time.Second * time.Duration(attempt)
			log.Printf("Retry attempt %d/%d after %v", attempt, c.config.MaxRetries, backoff)
			time.Sleep(backoff)
		}

		err := c.processMessage(ctx, msg)
		if err == nil {
			return nil
		}

		lastErr = err
		log.Printf("Processing failed (attempt %d/%d): %v", attempt+1, c.config.MaxRetries+1, err)
	}

	return fmt.Errorf("failed after %d retries: %w", c.config.MaxRetries, lastErr)
}

func (c *ConsumerImpl) processMessage(ctx context.Context, msg *kafkago.Message) error {
	processCtx, cancel := context.WithTimeout(ctx, c.config.ProcessingTimeout)
	defer cancel()

	event, err := FromJSON(msg.Value)
	if err != nil {
		return fmt.Errorf("failed to deserialize event: %w", err)
	}

	handlers := c.registry.GetHandlers(event.Type)
	if len(handlers) == 0 {
		log.Printf("No handlers registered for event type: %s", event.Type)
		return nil
	}

	var handlerErrors []error
	for idx, handler := range handlers {
		if err := processCtx.Err(); err != nil {
			return fmt.Errorf("context cancelled before handler %d: %w", idx, err)
		}

		if err := handler(event); err != nil {
			handlerErrors = append(handlerErrors, err)
			log.Printf("Handler %d error for event %s (type=%s): %v",
				idx, event.ID, event.Type, err)
		}
	}

	if len(handlerErrors) > 0 {
		return fmt.Errorf("handlers failed: %v", handlerErrors)
	}
	return nil
}

func (c *ConsumerImpl) Stop() error {
	c.mu.Lock()
	if !c.running {
		c.mu.Unlock()
		return nil
	}
	c.running = false
	c.mu.Unlock()

	if err := c.kafka.Consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}
	return nil
}

func (c *ConsumerImpl) Pause(partitions []kafkago.TopicPartition) error {
	return c.kafka.Consumer.Pause(partitions)
}

func (c *ConsumerImpl) Resume(partitions []kafkago.TopicPartition) error {
	return c.kafka.Consumer.Resume(partitions)
}

func (c *ConsumerImpl) Seek(partition kafkago.TopicPartition, timeoutMs int) error {
	return c.kafka.Consumer.Seek(partition, timeoutMs)
}

func (c *ConsumerImpl) CommitOffsets(offsets []kafkago.TopicPartition) error {
	_, err := c.kafka.Consumer.CommitOffsets(offsets)
	return err
}

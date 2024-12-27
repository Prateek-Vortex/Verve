package event

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
	Brokers = strings.Split(os.Getenv("event_broker"), ",") //[]string{"localhost:9092"}
	GroupID = os.Getenv("event_group")                      // "verve-group"
)

type Event interface {
	Publish(ctx context.Context, topic string, message interface{}) error
	Subscribe(ctx context.Context, topic string, handler func([]byte) error) error
	Close() error
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

type kafkaEvent struct {
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	logger   *slog.Logger
}

func NewKafkaEvent(logger *slog.Logger) (Event, error) {
	config := KafkaConfig{
		Brokers: Brokers,
		GroupID: GroupID,
	}
	// Producer config
	producerConfig := sarama.NewConfig()
	producerConfig.Producer.RequiredAcks = sarama.WaitForAll
	producerConfig.Producer.Retry.Max = 5
	producerConfig.Producer.Return.Successes = true

	// Create producer
	producer, err := sarama.NewSyncProducer(config.Brokers, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create producer: %w", err)
	}

	// Consumer config
	consumerConfig := sarama.NewConfig()
	consumerConfig.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	// Create consumer group
	consumer, err := sarama.NewConsumerGroup(config.Brokers, config.GroupID, consumerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &kafkaEvent{
		producer: producer,
		consumer: consumer,
		logger:   logger,
	}, nil
}

func (k *kafkaEvent) Publish(ctx context.Context, topic string, message interface{}) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic:     topic,
		Value:     sarama.ByteEncoder(bytes),
		Timestamp: time.Now(),
	}

	_, _, err = k.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (k *kafkaEvent) Subscribe(ctx context.Context, topic string, handler func([]byte) error) error {
	go func() {
		for {
			if err := k.consumer.Consume(ctx, []string{topic}, &consumerHandler{
				handler: handler,
				logger:  k.logger,
			}); err != nil {
				k.logger.Error("Error from consumer", "error", err)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

func (k *kafkaEvent) Close() error {
	if err := k.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}
	if err := k.consumer.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}
	return nil
}

type consumerHandler struct {
	handler func([]byte) error
	logger  *slog.Logger
}

func (h *consumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *consumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if err := h.handler(msg.Value); err != nil {
			h.logger.Error("Failed to process message", "error", err)
			continue
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

package consumer

import (
	"context"
	"log"

	"github.com/go-ecommerce-application/libs/kafka/config"
	"github.com/segmentio/kafka-go"
)

type MessageHandler func(ctx context.Context, message []byte) error

type Consumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

func NewConsumer(cfg *config.KafkaConfig, topic string, handler MessageHandler) (*Consumer, error) {
	if err := cfg.ValidateBrokers(); err != nil {
		return nil, err
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        cfg.Brokers,
		Topic:          topic,
		GroupID:        cfg.GroupID,
		StartOffset:    kafka.LastOffset,
		CommitInterval: 1 * 1000 * 1000 * 1000, // 1 second
	})

	return &Consumer{
		reader:  reader,
		handler: handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			message, err := c.reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("error reading message: %v", err)
				continue
			}

			if err := c.handler(ctx, message.Value); err != nil {
				log.Printf("error handling message: %v", err)
				continue
			}
		}
	}
}

func (c *Consumer) Close() error {
	if c.reader != nil {
		return c.reader.Close()
	}
	return nil
}

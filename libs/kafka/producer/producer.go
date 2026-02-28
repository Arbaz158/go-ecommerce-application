package producer

import (
	"context"
	"fmt"
	"log"

	"github.com/go-ecommerce-application/libs/kafka/config"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(cfg *config.KafkaConfig) (*Producer, error) {
	if err := cfg.ValidateBrokers(); err != nil {
		return nil, err
	}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{writer: writer}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, message []byte) error {
	err := p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: message,
		Topic: topic,
	})

	if err != nil {
		log.Printf("failed to publish message to topic %s: %v", topic, err)
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	if p.writer != nil {
		return p.writer.Close()
	}
	return nil
}

package config

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

func NewKafkaConfig(brokers []string, groupID string) *KafkaConfig {
	return &KafkaConfig{
		Brokers: brokers,
		GroupID: groupID,
	}
}

func (kc *KafkaConfig) GetDialer() *kafka.Dialer {
	return &kafka.Dialer{
		Timeout:   10000000000, // 10 seconds
		DualStack: true,
	}
}

func (kc *KafkaConfig) ValidateBrokers() error {
	if len(kc.Brokers) == 0 {
		return fmt.Errorf("no kafka brokers configured")
	}
	return nil
}

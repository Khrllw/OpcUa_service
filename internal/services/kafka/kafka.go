package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/interfaces"
	"strings"
)

type KafkaManager struct {
	producer *kafka.Producer
	cfg      *config.KafkaConfig
}

func NewKafkaManager(cfg *config.Config) (interfaces.KafkaManagerService, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": strings.Join(cfg.Kafka.Brokers, ","),
	})
	if err != nil {
		return nil, err
	}

	return &KafkaManager{
		producer: producer,
		cfg:      &cfg.Kafka,
	}, nil
}

func (k *KafkaManager) SendData(topic string, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	err = k.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          bytes,
	}, nil)
	if err != nil {
		return fmt.Errorf("produce error: %w", err)
	}

	return nil
}

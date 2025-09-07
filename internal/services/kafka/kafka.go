package kafka

import (
	"context"

	"opc_ua_service/internal/config"
	"opc_ua_service/internal/interfaces"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer создает новый экземпляр продюсера Kafka
func NewKafkaProducer(cfg *config.AppConfig) (interfaces.KafkaService, error) {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.KafkaBrokers...),
		Topic:    cfg.Kafka.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{writer: writer}, nil
}

// Produce отправляет сообщение в Kafka
func (p *KafkaProducer) Produce(ctx context.Context, key, value []byte) error {
	return p.writer.WriteMessages(ctx,
		kafka.Message{
			Key:   key,
			Value: value,
		},
	)
}

// Close закрывает соединение с Kafka
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

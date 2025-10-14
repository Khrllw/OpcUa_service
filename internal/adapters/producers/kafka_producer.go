package producers

import (
	"context"
	"github.com/segmentio/kafka-go"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"
	"time"
)

type KafkaProducer struct {
	writer         *kafka.Writer
	logger         *logging.Logger
	usecase        interfaces.Usecases
	publishTimeout time.Duration
	stopChan       chan struct{}
	retryDelay     time.Duration
	batchSize      int
}

// NewKafkaProducer создает новый экземпляр продюсера Kafka
func NewKafkaProducer(cfg *config.Config, usecase interfaces.Usecases, parentLogger *logging.Logger) (interfaces.KafkaProducer, error) {
	kafkaLogger := parentLogger.WithPrefix("KAFKA PRODUCER")
	kafkaLogger.Info("Kafka initialized",
		"component", "GENERAL",
	)

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.App.Kafka.KafkaBrokers...),
		Topic:    cfg.App.Kafka.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	return &KafkaProducer{
		writer:         writer,
		usecase:        usecase,
		logger:         kafkaLogger,
		publishTimeout: time.Duration(cfg.App.Kafka.KafkaTimeout) * time.Second,
		retryDelay:     time.Duration(cfg.App.Kafka.KafkaRetryDelay) * time.Second,
		batchSize:      cfg.App.Kafka.KafkaBatchSize,
		stopChan:       make(chan struct{})}, nil
}

func (p *KafkaProducer) Start(ctx context.Context) {
	go p.publisherWorker()
	go p.cleanerWorker()
}

// Stop останавливает воркеры
func (p *KafkaProducer) Stop() {
	close(p.stopChan)
}

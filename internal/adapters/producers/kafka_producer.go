package producers

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/internal/middleware/logging"
	"strconv"
	"time"
)

type KafkaProducer struct {
	writer         *kafka.Writer
	logger         *logging.Logger
	usecase        interfaces.Usecases
	publishTimeout time.Duration
	stopChan       chan struct{}
	retryDelay     time.Duration
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
		stopChan:       make(chan struct{})}, nil
}

// Start запускает фонового воркера, который читает данные из БД и отправляет в Kafka
func (p *KafkaProducer) Start(ctx context.Context) {
	p.logger.Info("Starting KafkaProducer background worker...")

	go func() {
		ticker := time.NewTicker(p.publishTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				p.publishAll(ctx)
			case <-p.stopChan:
				p.logger.Info("KafkaProducer stopped")
				return
			}
		}
	}()
}

// publishAll читает данные из БД и пытается отправить в Kafka с retry
func (p *KafkaProducer) publishAll(ctx context.Context) {
	records, err := p.usecase.GetAllPollData()
	if err != nil {
		p.logger.Warn("Failed to get poll data from DB", "error", err)
		return
	}

	for _, record := range records {
		if err := p.publishRecordWithRetry(record); err != nil {
			p.logger.Warn("Failed to publish record after retries", "id", record.ID, "error", err)
		}
	}
}

func (p *KafkaProducer) publishRecordWithRetry(record entities.PollData) error {
	payload, err := json.Marshal(record)
	if err != nil {
		p.logger.Warn("Failed to marshal payload", "id", record.ID, "error", err)
		return err
	}

	key := []byte(strconv.FormatUint(uint64(record.ID), 10))

	for {
		select {
		case <-p.stopChan:
			p.logger.Info("KafkaProducer stopped")
			return nil
		default:
			// локальный контекст только на одну попытку
			timeoutCtx, cancel := context.WithTimeout(context.Background(), p.publishTimeout)
			err := p.writer.WriteMessages(timeoutCtx, kafka.Message{
				Key:   key,
				Value: payload,
			})
			cancel()

			if err == nil {
				// успешно отправлено — удаляем из БД
				if err := p.usecase.DeletePollDataByID(record.ID); err != nil {
					p.logger.Warn("Failed to delete poll data from DB", "id", record.ID, "error", err)
				}
				return nil
			}

			p.logger.Warn("Kafka produce failed, retrying...", "id", record.ID, "error", err)
			p.reconnectIfNeeded()
			time.Sleep(p.retryDelay)
		}
	}
}

// reconnectIfNeeded пытается "восстановить" соединение с Kafka после ошибки
func (p *KafkaProducer) reconnectIfNeeded() {
	// закрываем старый writer и создаём новый
	if err := p.writer.Close(); err != nil {
		p.logger.Warn("Failed to close Kafka writer during reconnect", "error", err)
	}

	p.logger.Info("Reconnecting Kafka writer...")
	// создаём новый writer с теми же параметрами
	p.writer = &kafka.Writer{
		Addr:     p.writer.Addr,
		Topic:    p.writer.Topic,
		Balancer: &kafka.LeastBytes{},
	}
	p.logger.Info("Kafka writer reconnected")
}

// Close закрывает соединение с Kafka
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}

// Stop останавливает фонового воркера
func (p *KafkaProducer) Stop() {
	err := p.writer.Close()
	if err != nil {
		p.logger.Warn("Failed to close Kafka writer", "error", err)
	}
	close(p.stopChan)
}

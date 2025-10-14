package producers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"opc_ua_service/pkg/errors"
	"strconv"
	"time"
)

// publisherWorker читает БД пакетами и публикует в Kafka
func (p *KafkaProducer) publisherWorker() {
	p.logger.Info("Publisher worker started")
	ticker := time.NewTicker(p.publishTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.publishAll()
		case <-p.stopChan:
			p.logger.Info("Publisher worker stopped")
			return
		}
	}
}

// publishAll читает данные из БД и отправляет batch в Kafka с retry
func (p *KafkaProducer) publishAll() {
	for {
		records, err := p.usecase.GetPollDataBatch(p.batchSize)
		if err != nil {
			p.logger.Warn("Failed to fetch poll data batch", "error", err)
			return
		}
		if len(records) == 0 {
			// данных нет — выходим
			return
		}

		// формируем batch сообщений
		messages := make([]kafka.Message, 0, len(records))
		ids := make([]uint, 0, len(records)) // сразу собираем ID
		for _, r := range records {
			payload, _ := json.Marshal(r)
			messages = append(messages, kafka.Message{
				Key:   []byte(strconv.FormatUint(uint64(r.ID), 10)),
				Value: payload,
			})
			ids = append(ids, r.ID)
		}

	retryLoop: // ← метка, чтобы управлять внутренним циклом
		for {
			select {
			case <-p.stopChan:
				p.logger.Info("Publisher worker stopping...")
				return
			default:
				ctx, cancel := context.WithTimeout(context.Background(), p.publishTimeout)
				err := p.writer.WriteMessages(ctx, messages...)
				cancel()

				if err == nil {
					// успешно отправлено — помечаем processed
					if err := p.usecase.MarkPollDataProcessed(ids); err != nil && !errors.Is(err, errors.ErrEmptyAction) {
						p.logger.Warn("Failed to mark records as processed", "error", err)
					}
					// выходим из retryLoop и идём к следующему batch
					break retryLoop
				}

				p.logger.Warn("Kafka publish failed, retrying...", "error", err)
				if err := p.reconnectWriter(); err != nil {
					p.logger.Warn("Failed to reconnect writer", "error", err)
				}
				time.Sleep(p.retryDelay)
			}
		}
	}
}

func (p *KafkaProducer) reconnectWriter() error {
	p.logger.Info("Reconnecting Kafka writer...")

	// Закрываем старый writer, если он существует
	if p.writer != nil {
		if err := p.writer.Close(); err != nil {
			p.logger.Warn("Failed to close old Kafka writer", "error", err)
		}
	}

	// Попытка создать новый writer с retry
	for {
		select {
		case <-p.stopChan:
			p.logger.Info("KafkaProducer stopped, aborting reconnect")
			return fmt.Errorf("producer stopped")
		default:
			writer := &kafka.Writer{
				Addr:     p.writer.Addr, // используем те же брокеры
				Topic:    p.writer.Topic,
				Balancer: &kafka.LeastBytes{},
			}

			// Тестовая попытка записи пустого сообщения для проверки подключения
			testCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err := writer.WriteMessages(testCtx)
			cancel()

			if err == nil {
				p.writer = writer
				p.logger.Info("Kafka writer reconnected")
				return nil
			}

			p.logger.Warn("Kafka reconnect failed, retrying...", "error", err)
			time.Sleep(p.retryDelay)
		}
	}
}

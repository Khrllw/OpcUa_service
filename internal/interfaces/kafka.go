package interfaces

import (
	"context"
)

// KafkaProducer определяет контракт для отправки данных во внешние системы
type KafkaProducer interface {
	Start(ctx context.Context)
	Stop()
	Close() error
}

package kafka

import (
	"opc_ua_service/internal/config"
	"opc_ua_service/internal/interfaces"
)

type Kafka struct {
	Manager interfaces.KafkaManagerService
}

func NewKafka(cfg *config.Config) (*Kafka, error) {
	manager, err := NewKafkaManager(cfg)
	if err != nil {
		return nil, err
	}

	return &Kafka{
		Manager: manager,
	}, nil
}

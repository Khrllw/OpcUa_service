package interfaces

type Kafka interface {
	KafkaManagerService
}

type KafkaManagerService interface {
	SendData(topic string, data interface{}) error
}

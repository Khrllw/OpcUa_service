package models

import (
	"time"
)

// AnonymousConnection конфигурация подключения OPC UA анонимно
type AnonymousConnection struct {
	EndpointURL  string
	Timeout      time.Duration
	Manufacturer string
	Model        string
}

func (*AnonymousConnection) GetType() ConnectionTypeEnum {
	return ConnectionAnonymous
}

func (ac *AnonymousConnection) GetManufacturer() string {
	return ac.Model
}

func (ac *AnonymousConnection) GetModel() string {
	return ac.Model
}

func (ac *AnonymousConnection) GetEndpointURL() string {
	return ac.EndpointURL
}

func (ac *AnonymousConnection) GetTimeout() time.Duration {
	return ac.Timeout
}

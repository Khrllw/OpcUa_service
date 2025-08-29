package models

import (
	"time"
)

// PasswordConnection конфигурация подключения OPC UA по паролю
type PasswordConnection struct {
	EndpointURL  string
	Username     string
	Password     string
	Policy       string
	Mode         string
	Timeout      time.Duration
	Manufacturer string
	Model        string
}

func (*PasswordConnection) GetType() ConnectionTypeEnum {
	return ConnectionPassword
}

func (pc *PasswordConnection) GetManufacturer() string {
	return pc.Model
}

func (pc *PasswordConnection) GetModel() string {
	return pc.Model
}

func (pc *PasswordConnection) GetEndpointURL() string {
	return pc.EndpointURL
}

func (pc *PasswordConnection) GetTimeout() time.Duration {
	return pc.Timeout
}

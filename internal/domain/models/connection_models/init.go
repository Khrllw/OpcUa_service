package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/url"
	"time"
)

type ConnectionConfig struct {
	//Type   ConnectionTypeEnum
	Config ConnectionConfigImpl
}

type ConnectionConfigImpl interface {
	GetType() ConnectionTypeEnum
	GetManufacturer() string
	GetModel() string
	GetEndpointURL() string
	GetTimeout() time.Duration
}

// ---------------------------------------------------------------------------------------------------------------

// ConnectionTypeEnum - допустимые типы аутентификации
type ConnectionTypeEnum string

const (
	ConnectionAnonymous   ConnectionTypeEnum = "anonymous"
	ConnectionPassword    ConnectionTypeEnum = "password"
	ConnectionCertificate ConnectionTypeEnum = "certificate"
)

// Validate для проверки enum ConnectionTypeEnum
func (a ConnectionTypeEnum) Validate() error {
	switch a {
	case ConnectionAnonymous, ConnectionPassword, ConnectionCertificate:
		return nil
	default:
		return validator.ValidationErrors{}
	}
}

// ---------------------------------------------------------------------------------------------------------------

// ConnectionStatusEnum - статусы подключения к машине
type ConnectionStatusEnum string

const (
	ConnectionStatusConnected    ConnectionStatusEnum = "connected"
	ConnectionStatusDisconnected ConnectionStatusEnum = "disconnected"
	ConnectionStatusPolled       ConnectionStatusEnum = "polled"
)

// ---------------------------------------------------------------------------------------------------------------

func (c ConnectionConfig) MarshalJSON() ([]byte, error) {
	type wrapper struct {
		Type   string      `json:"type"`
		Config interface{} `json:"connection"`
	}

	return json.Marshal(wrapper{
		Type:   string(c.Config.GetType()),
		Config: c.Config,
	})
}

func (c *ConnectionConfig) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type   string          `json:"type"`
		Config json.RawMessage `json:"connection"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch raw.Type {
	case "password":
		var cfg PasswordConnection
		if err := json.Unmarshal(raw.Config, &cfg); err != nil {
			return err
		}
		c.Config = &cfg

	case "certificate":
		var cfg CertificateConnection
		if err := json.Unmarshal(raw.Config, &cfg); err != nil {
			return err
		}
		c.Config = &cfg

	case "anonymous":
		var cfg AnonymousConnection
		if err := json.Unmarshal(raw.Config, &cfg); err != nil {
			return err
		}
		c.Config = &cfg

	default:
		return fmt.Errorf("unknown config type: %s", raw.Type)
	}

	return nil
}

func (c ConnectionConfig) EndpointURLHostPort() string {
	u, _ := url.Parse(c.Config.GetEndpointURL())
	return fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
}

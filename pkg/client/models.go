package client

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// ConnectionRequest - данные для подключения клиента ЧПУ
type ConnectionRequest struct {
	ConnectionType ConnectionTypeEnum      `json:"connectionType,omitempty" binding:"required" example:"password"`
	Username       string                  `json:"username,omitempty" example:"client1"`         // для password
	Password       string                  `json:"password,omitempty" example:"secret"`          // для password
	Certificate    string                  `json:"certificate,omitempty" example:"cert-abc-123"` // для certificate
	Key            string                  `json:"key,omitempty" example:"secret"`
	EndpointURL    string                  `json:"endpointURL" example:"opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC"`
	Policy         SecurityPolicyEnum      `json:"policy,omitempty" example:"Basic256Sha256"` // OPC UA SecurityPolicy
	Mode           MessageSecurityModeEnum `json:"mode,omitempty" example:"SignAndEncrypt"`   // OPC UA MessageSecurityMode
	Timeout        int                     `json:"timeout,omitempty" example:"30"`
	Manufacturer   string                  `json:"manufacturer" binding:"required" example:"Heidenhain"`
	Model          string                  `json:"model" binding:"required" example:"TNC640"`
}

// DisconnectRequest - отключение станка
type DisconnectRequest struct {
	UUID string `json:"ID" binding:"required"`
}

type IDRequest struct {
	ID uint `json:"ID" binding:"required"`
}

type StartPollingRequest struct {
	ID      uint `json:"ID" binding:"required"`
	Timeout int  `json:"timeout" binding:"required"`
}

// CheckConnectionRequest - проверка соединения
type CheckConnectionRequest struct {
	ID uint `json:"ID" binding:"required"`
}

// GetControlProgramRequest - получение управляющей программы
type GetControlProgramRequest struct {
	ID uint `json:"ID" binding:"required"`
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

// SecurityPolicyEnum - допустимые политики безопасности OPC UA
type SecurityPolicyEnum string

const (
	PolicyNone           SecurityPolicyEnum = "None"
	PolicyBasic128Rsa15  SecurityPolicyEnum = "Basic128Rsa15"
	PolicyBasic256       SecurityPolicyEnum = "Basic256"
	PolicyBasic256Sha256 SecurityPolicyEnum = "Basic256Sha256"
)

// Validate проверяет корректность SecurityPolicyEnum
func (p SecurityPolicyEnum) Validate() error {
	switch p {
	case PolicyNone, PolicyBasic128Rsa15, PolicyBasic256, PolicyBasic256Sha256:
		return nil
	default:
		return fmt.Errorf("invalid security policy: %s", p)
	}
}

// ---------------------------------------------------------------------------------------------------------------

// MessageSecurityModeEnum - допустимые режимы шифрования сообщений OPC UA
type MessageSecurityModeEnum string

const (
	ModeNone           MessageSecurityModeEnum = "None"
	ModeSign           MessageSecurityModeEnum = "Sign"
	ModeSignAndEncrypt MessageSecurityModeEnum = "SignAndEncrypt"
)

// Validate проверяет корректность MessageSecurityModeEnum
func (m MessageSecurityModeEnum) Validate() error {
	switch m {
	case ModeNone, ModeSign, ModeSignAndEncrypt:
		return nil
	default:
		return fmt.Errorf("invalid message security mode: %s", m)
	}
}

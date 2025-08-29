package models

import (
	"fmt"
	"net/http"
)

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

// ---------------------------------------------------------------------------------------------------------------

// ConnectionStatusEnum - детальные статусы соединения OPC UA
type ConnectionStatusEnum string

const (
	// Базовые статусы
	StatusOK       ConnectionStatusEnum = "OK"
	StatusFail     ConnectionStatusEnum = "FAIL"
	StatusNotFound ConnectionStatusEnum = "NOT FOUND"

	// Детальные статусы соединения
	StatusHealthy      ConnectionStatusEnum = "HEALTHY"      // Соединение здорово и работает
	StatusUnhealthy    ConnectionStatusEnum = "UNHEALTHY"    // Соединение есть, но есть проблемы
	StatusConnected    ConnectionStatusEnum = "CONNECTED"    // Соединение установлено
	StatusDisconnected ConnectionStatusEnum = "DISCONNECTED" // Соединение разорвано
	StatusError        ConnectionStatusEnum = "ERROR"        // Ошибка соединения
	StatusTimeout      ConnectionStatusEnum = "TIMEOUT"      // Таймаут соединения
	StatusAuthFailed   ConnectionStatusEnum = "AUTH_FAILED"  // Ошибка аутентификации
)

// Validate для проверки enum ConnectionStatusEnum
func (s ConnectionStatusEnum) Validate() error {
	switch s {
	case StatusOK, StatusFail,
		StatusHealthy, StatusUnhealthy,
		StatusConnected, StatusDisconnected,
		StatusError, StatusTimeout, StatusAuthFailed, StatusNotFound:
		return nil
	default:
		return fmt.Errorf("invalid connection status: %s", s)
	}
}

// String возвращает строковое представление статуса
func (s ConnectionStatusEnum) String() string {
	return string(s)
}

// IsOK проверяет, является ли статус успешным
func (s ConnectionStatusEnum) IsOK() bool {
	switch s {
	case StatusOK, StatusHealthy, StatusConnected:
		return true
	default:
		return false
	}
}

// IsError проверяет, является ли статус ошибкой
func (s ConnectionStatusEnum) IsError() bool {
	switch s {
	case StatusFail, StatusError, StatusTimeout,
		StatusAuthFailed, StatusNotFound:
		return true
	default:
		return false
	}
}

// IsActive проверяет, активно ли соединение
func (s ConnectionStatusEnum) IsActive() bool {
	switch s {
	case StatusHealthy, StatusConnected:
		return true
	default:
		return false
	}
}

// ToHTTPStatus преобразует статус соединения в HTTP статус
func (s ConnectionStatusEnum) ToHTTPStatus() int {
	switch s {
	case StatusOK, StatusHealthy, StatusConnected:
		return http.StatusOK
	case StatusUnhealthy:
		return http.StatusOK // 200, но с предупреждением
	case StatusDisconnected:
		return http.StatusGone
	case StatusNotFound:
		return http.StatusNotFound
	case StatusError, StatusFail, StatusTimeout,
		StatusAuthFailed:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// GetDescription возвращает человеко-читаемое описание статуса
func (s ConnectionStatusEnum) GetDescription() string {
	switch s {
	case StatusOK:
		return "Operation completed successfully"
	case StatusFail:
		return "Operation failed"
	case StatusHealthy:
		return "Connection is healthy and responsive"
	case StatusUnhealthy:
		return "Connection is established but experiencing issues"
	case StatusConnected:
		return "Connection successfully established"
	case StatusDisconnected:
		return "Connection has been disconnected"
	case StatusError:
		return "Connection error occurred"
	case StatusTimeout:
		return "Connection timeout occurred"
	case StatusAuthFailed:
		return "Authentication failed"
	case StatusNotFound:
		return "Connection not found"
	default:
		return "Unknown connection status"
	}
}

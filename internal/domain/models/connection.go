package models

import (
	"github.com/google/uuid"
	models "opc_ua_service/internal/domain/models/connection_models"
	"time"
)

// ConnectionRequest - данные для аутентификации клиента ЧПУ
type ConnectionRequest struct {
	ConnectionType models.ConnectionTypeEnum `json:"connectionType,omitempty" binding:"required" example:"password"`
	Username       string                    `json:"username,omitempty" example:"client1"`         // для password
	Password       string                    `json:"password,omitempty" example:"secret"`          // для password
	Certificate    string                    `json:"certificate,omitempty" example:"cert-abc-123"` // для certificate
	Key            string                    `json:"key,omitempty" example:"secret"`
	EndpointURL    string                    `json:"endpointURL" example:"opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC"`
	Policy         SecurityPolicyEnum        `json:"policy,omitempty" example:"Basic256Sha256"` // OPC UA SecurityPolicy
	Mode           MessageSecurityModeEnum   `json:"mode,omitempty" example:"SignAndEncrypt"`   // OPC UA MessageSecurityMode
	Timeout        int                       `json:"timeout,omitempty" example:"30"`
	Manufacturer   string                    `json:"manufacturer" binding:"required" example:"Heidenhain"`
	Model          string                    `json:"model" binding:"required" example:"TNC640"`
}

// DisconnectRequest - отключение станка
type DisconnectRequest struct {
	UUID string `json:"UUID" binding:"required"`
}

type UUIDRequest struct {
	UUID string `json:"UUID" binding:"required"`
}

// CheckConnectionRequest - проверка соединения
type CheckConnectionRequest struct {
	UUID string `json:"UUID" binding:"required"`
}

// GetControlProgramRequest - получение управляющей программы
type GetControlProgramRequest struct {
	UUID string `json:"UUID" binding:"required"`
}

// ---------------------------------------------------------------------------------------------------------------

// UUIDResponse - UUID станка
type UUIDResponse struct {
	UUID string `json:"UUID" binding:"required"`
}

// DisconnectResponse - успешное отключение станка
type DisconnectResponse struct {
	Disconnected bool `json:"disconnected" binding:"required"`
}

// CheckConnectionResponse - успешное отключение станка
type CheckConnectionResponse struct {
	Connected bool `json:"connected" binding:"required"`
}

// ConnectionResponse - ответ при успешной аутентификации
type ConnectionResponse struct {
	Status         ConnectionStatusEnum    `json:"status" binding:"required" example:"OK"`
	SessionID      string                  `json:"sessionID"`
	ConnectionInfo *ConnectionInfoResponse `json:"connectionInfo"`
}

// ConnectionInfoResponse - ответ с информацией о подключении
type ConnectionInfoResponse struct {
	Status      ConnectionStatusEnum    `json:"status"`
	Description string                  `json:"description"`
	UUID        uuid.UUID               `json:"UUID"`
	SessionID   string                  `json:"sessionID" example:"ns=3;i=3093118269"`
	Config      models.ConnectionConfig `json:"config"`
	CreatedAt   time.Time               `json:"createdAt" example:"2025-08-22T12:00:00Z"`
	LastUsed    time.Time               `json:"lastUsed" example:"2025-08-22T12:05:00Z"`
	UseCount    int64                   `json:"useCount" example:"1"`
}

// CheckConnectionWithInfoResponse - ответ проверки соединения
type CheckConnectionWithInfoResponse struct {
	Status         ConnectionStatusEnum    `json:"status"`
	SessionID      string                  `json:"sessionID"`
	ConnectionInfo *ConnectionInfoResponse `json:"connectionInfo"`
}

// ConnectionPoolResponse - информация об активных подключениях
type ConnectionPoolResponse struct {
	PoolSize    int                       `json:"poolSize"`
	Connections []*ConnectionInfoResponse `json:"connections"`
}

type PollingResponse struct {
	Polled bool `json:"polled" binding:"required"`
}

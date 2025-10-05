package connection

import (
	"github.com/google/uuid"
	models "opc_ua_service/internal/domain/models/connection_types"
	"time"
)

// ConnectionRequest - данные для подключения клиента ЧПУ
type ConnectionRequest struct {
	ConnectionType models.ConnectionTypeEnum `json:"connectionType,omitempty" binding:"required" example:"password"`
	Username       string                    `json:"username,omitempty" example:"client1"`         // для password
	Password       string                    `json:"password,omitempty" example:"secret"`          // для password
	Certificate    string                    `json:"certificate,omitempty" example:"cert-abc-123"` // для certificate
	Key            string                    `json:"key,omitempty" example:"secret"`
	EndpointURL    string                    `json:"endpointURL" example:"opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC"`
	Policy         SecurityPolicyEnum        `json:"policy,omitempty" example:"Basic256Sha256"` // OPC UA SecurityPolicy
	Mode           MessageSecurityModeEnum   `json:"mode,omitempty" example:"SignAndEncrypt"`   // OPC UA MessageSecurityMode
	Manufacturer   string                    `json:"manufacturer" binding:"required" example:"Heidenhain"`
	Model          string                    `json:"model" binding:"required" example:"TNC640"`
}

// DisconnectRequest - отключение станка
type DisconnectRequest struct {
	ID uint `json:"ID" binding:"required"`
}

// StartPollingRequest - начало опроса станка
type StartPollingRequest struct {
	ID      uint `json:"ID" binding:"required"`
	Timeout int  `json:"timeout" binding:"required"`
}

type IDRequest struct {
	ID uint `json:"ID" binding:"required"`
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

// IDResponse - ID станка
type IDResponse struct {
	ID uint `json:"ID" binding:"required"`
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
	MachineID   uint64                  `json:"machineID" example:"1"`
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

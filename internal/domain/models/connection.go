package models

import (
	"fmt"
	"github.com/awcullen/opcua/ua"
	"github.com/go-playground/validator/v10"
	"time"
)

// ConnectionTypeEnum - допустимые типы аутентификации
type ConnectionTypeEnum string

const (
	ConnectionAnonymous   ConnectionTypeEnum = "anonymous"
	ConnectionPassword    ConnectionTypeEnum = "password"
	ConnectionCertificate ConnectionTypeEnum = "certificate"
)

// Validate для проверки enum AuthTypeEnum
func (a ConnectionTypeEnum) Validate() error {
	switch a {
	case ConnectionAnonymous, ConnectionPassword, ConnectionCertificate:
		return nil
	default:
		return validator.ValidationErrors{}
	}
}

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

// ClientLoginRequest - данные для аутентификации клиента ЧПУ
type ConnectionRequest struct {
	ConnectionType ConnectionTypeEnum      `json:"connectionType,omitempty" example:"password"`
	Username       string                  `json:"username,omitempty" example:"client1"`         // для password
	Password       string                  `json:"password,omitempty" example:"secret"`          // для password
	Certificate    string                  `json:"certificate,omitempty" example:"cert-abc-123"` // для certificate
	Key            string                  `json:"key,omitempty" example:"secret"`
	EndpointURL    string                  `json:"endpointURL" example:"opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC"`
	Policy         SecurityPolicyEnum      `json:"policy,omitempty" example:"Basic256Sha256"` // OPC UA SecurityPolicy
	Mode           MessageSecurityModeEnum `json:"mode,omitempty" example:"SignAndEncrypt"`   // OPC UA MessageSecurityMode
	Timeout        int                     `json:"timeout,omitempty" example:"30"`            // таймаут в секундах
}

// ---------------------------------------------------------------------------------------------------------------

// AuthStatusEnum - допустимые статусы ответа
type AuthStatusEnum string

const (
	StatusOK   AuthStatusEnum = "OK"
	StatusFail AuthStatusEnum = "FAIL"
)

// Validate для проверки enum AuthStatusEnum
func (s AuthStatusEnum) Validate() error {
	switch s {
	case StatusOK, StatusFail:
		return nil
	default:
		return validator.ValidationErrors{}
	}
}

// ConnectionAuthResponse - ответ при успешной аутентификации
type ConnectionAuthResponse struct {
	Status         AuthStatusEnum          `json:"status" binding:"required,oneof=OK FAIL" example:"OK"` // Статус ответа
	Token          string                  `json:"token"`                                                // Токен установленной сессии
	ConnectionInfo *ConnectionInfoResponse `json:"connectionInfo"`
}

type ConnectionInfoResponse struct {
	SessionID string           `json:"sessionId" example:"ns=3;i=3093118269"`
	Config    ConnectionConfig `json:"config"`
	CreatedAt time.Time        `json:"createdAt" example:"2025-08-22T12:00:00Z"`
	LastUsed  time.Time        `json:"lastUsed" example:"2025-08-22T12:05:00Z"`
	UseCount  int64            `json:"useCount" example:"1"`
	IsHealthy bool             `json:"isHealthy" example:"true"`
}

func ToResponse(info *ConnectionInfo) ConnectionInfoResponse {
	if info == nil {
		return ConnectionInfoResponse{}
	}

	return ConnectionInfoResponse{
		SessionID: info.SessionID,
		Config:    info.Config,
		CreatedAt: info.CreatedAt,
		LastUsed:  info.LastUsed,
		UseCount:  info.UseCount,
		IsHealthy: info.IsHealthy,
	}
}

type DisconnectRequest struct {
	SessionID string `json:"sessionID" binding:"required"`
}

type CheckConnectionRequest struct {
	SessionID string `json:"sessionID" binding:"required"`
}

func (ci *ConnectionInfo) GetRelevantNodeIDs() []ua.NodeIDNumeric {
	switch ci.Manufacturer {
	case "ACME", "Heidenhain":
		switch ci.Model {
		case "TNC640":
			// Здесь перечисляем NodeID, которые нам нужны для Heidenhain TNC640
			return []ua.NodeIDNumeric{

				ua.NewNodeIDNumeric(1, 100024), // OperatingMode

				// ------------------------ TOOL ------------------------
				ua.NewNodeIDNumeric(1, 100039), // CurrentToolName

				ua.NewNodeIDNumeric(1, 100003), // CutterLocation

				// ------------------------ FEED ------------------------
				ua.NewNodeIDNumeric(1, 100025), // FeedOverride
				ua.NewNodeIDNumeric(1, 100026), // FeedOverrideEURange
				ua.NewNodeIDNumeric(1, 300002), // FeedOverrideEngineeringUnits

				// ------------------------ RAPID ------------------------
				ua.NewNodeIDNumeric(1, 100029), // RapidOverride
				ua.NewNodeIDNumeric(1, 100030), // RapidOverrideEURange
				ua.NewNodeIDNumeric(1, 300004), // RapidOverrideEngineeringUnits
				ua.NewNodeIDNumeric(1, 100031), // RapidTraverseActive

				// ------------------------ SPEED ------------------------
				ua.NewNodeIDNumeric(1, 100027), // SpeedOverride
				ua.NewNodeIDNumeric(1, 100028), // SpeedOverrideEURange
				ua.NewNodeIDNumeric(1, 300003), // SpeedOverrideEngineeringUnits

				// ------------------------ TIME ------------------------
				ua.NewNodeIDNumeric(1, 56031), // ControlUpTime
				ua.NewNodeIDNumeric(1, 56033), // MachineUpTime
				ua.NewNodeIDNumeric(1, 56032), // ProgramExecutionTime

				// ---------------  PROGRAM  -------------------
				ua.NewNodeIDNumeric(1, 51002),  // CurrentState
				ua.NewNodeIDNumeric(1, 100005), // CurrentCall
				ua.NewNodeIDNumeric(1, 100006), // ExecutionStack
				ua.NewNodeIDNumeric(1, 100022), // ActiveProgramName

				// ----------------- EXECUTION STATE -------------------------
				ua.NewNodeIDNumeric(1, 100010), // ExecutionStateCurrentState
				ua.NewNodeIDNumeric(1, 100008), // ExecutionStateLastTransition
			}
		default:
			return nil
		}
	default:
		return nil
	}
}

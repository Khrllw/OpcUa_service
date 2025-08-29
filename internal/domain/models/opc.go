package models

import (
	"context"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/internal/domain/models/opc_custom"
	"sync"
	"time"
)

type ControlProgramInfoRequest struct {
	ExecutionStack []opc_custom.ProgramPositionDataType `json:"ExecutionStack"`
}

// ConnectorStats статистика менеджера подключений
type ConnectorStats struct {
	TotalConnections  int64
	FailedConnections int64
	ActiveConnections int64
	PoolSize          int64
}

// ConnectionInfo информация о подключении
type ConnectionInfo struct {
	Conn      *client.Client
	Ctx       context.Context
	Cancel    context.CancelFunc
	SessionID string
	Config    models.ConnectionConfig
	CreatedAt time.Time
	LastUsed  time.Time
	UseCount  int64
	IsHealthy bool
	Mu        sync.RWMutex

	Manufacturer string
	Model        string
}

func (ci *ConnectionInfo) GetRelevantNodeIDs() []ua.NodeIDNumeric {
	switch ci.Manufacturer {
	case "ACME", "Heidenhain":
		switch ci.Model {
		case "TNC640":
			// Здесь перечисляем NodeID, которые нам нужны для Heidenhain TNC640
			return []ua.NodeIDNumeric{
				ua.NewNodeIDNumeric(1, 56004), // SerialNumber

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

func (ci *ConnectionInfo) GetControlProgramNodes() []ua.NodeIDNumeric {
	switch ci.Manufacturer {
	case "ACME", "Heidenhain":
		switch ci.Model {
		case "TNC640":
			// Здесь перечисляем NodeID, которые нам нужны для Heidenhain TNC640
			return []ua.NodeIDNumeric{
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

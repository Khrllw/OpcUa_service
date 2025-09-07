package interfaces

import (
	"github.com/awcullen/opcua/ua"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/pkg/machine_models"
	"opc_ua_service/pkg/opc_custom"
)

// MachineData — общий интерфейс для всех моделей станков
type MachineData interface {
	ConvertNodeToMachineData(nodeID string, v any) error
	GetExecutionStack() ([]opc_custom.ProgramPositionDataType, error)
	ToJSON() string
	ToResponse() models.MachineDataResponse
	GetRelevantNodeIDs() []ua.NodeIDNumeric
	GetMachineID() (*string, error)
}

func MachineDataFactory(manufacturer, model string) MachineData {
	switch manufacturer {
	case "ACME", "Heidenhain":
		switch model {
		case "TNC640":
			return &machine_models.HeidenhainTNC640Data{}
		case "TNC620":
			return &machine_models.HeidenhainTNC640Data{}
		default:
			return nil // неизвестная модель для этого производителя
		}
	default:
		return nil // неизвестный производитель
	}
}

package interfaces

import (
	"github.com/awcullen/opcua/ua"
	"opc_ua_service/internal/domain/models/machine_models"
)

// MachineData — общий интерфейс для всех моделей станков
type MachineData interface {
	DecodeFromVariant(v ua.Variant) error // Декодировать из OPC UA значения
	ConvertNodeToMachineData(nodeID string, v any) error
}

func MachineDataFactory(manufacturer, model string) MachineData {
	switch manufacturer {
	case "ACME", "Heidenhain":
		switch model {
		case "TNC640":
			return &machine_models.HeidenhainTNC640Data{}
		default:
			return nil // неизвестная модель для этого производителя
		}
	default:
		return nil // неизвестный производитель
	}
}

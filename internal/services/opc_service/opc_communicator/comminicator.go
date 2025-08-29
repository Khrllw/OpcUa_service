package opc_communicator

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/ua"
	"github.com/google/uuid"
	"log"
	_ "opc_ua_service/internal/domain/models"
	"opc_ua_service/internal/domain/models/opc_custom"
	"opc_ua_service/internal/interfaces"
	"sync"
)

// OpcCommunicator содержит коннектор с пулом соединений
type OpcCommunicator struct {
	connector     interfaces.OpcConnectorService
	pollCancelMap map[uuid.UUID]context.CancelFunc
	mu            sync.Mutex
}

// NewOpcCommunicator создает новый экземпляр OpcCommunicator
func NewOpcCommunicator(connector interfaces.OpcConnectorService) interfaces.OpcCommunicatorService {
	return &OpcCommunicator{
		connector:     connector,
		pollCancelMap: make(map[uuid.UUID]context.CancelFunc),
	}
}

func (oc *OpcCommunicator) ReadMachineData(id uuid.UUID) (interfaces.MachineData, error) {
	connInfo, err := oc.connector.GetConnectionInfoByUUID(id)
	if err != nil {
		return nil, fmt.Errorf("connection not found: %w", err)
	}

	// Создаём объект данных машины через фабрику
	machine := interfaces.MachineDataFactory(connInfo.Manufacturer, connInfo.Model)
	if machine == nil {
		return nil, fmt.Errorf("unsupported machine type: %s %s", connInfo.Manufacturer, connInfo.Model)
	}

	// Список NodeID, которые нужно читать для этой машины
	nodeIDs := connInfo.GetRelevantNodeIDs()
	if len(nodeIDs) == 0 {
		return nil, fmt.Errorf("no nodes defined for machine %s %s", connInfo.Manufacturer, connInfo.Model)
	}

	// Считываем значение каждого узла
	for _, nodeID := range nodeIDs {
		val, err := oc.readNodeValue(connInfo.Ctx, connInfo.Conn, nodeID)
		if err != nil {
			log.Printf("Failed to read node %s: %v", nodeID, err)
			continue
		}

		// Декодируем значение в структуру
		if err := machine.ConvertNodeToMachineData(nodeID.String(), val); err != nil {
			log.Printf("Failed to convert node %s: %v", nodeID, err)
			continue
		}
	}

	return machine, nil
}

func (oc *OpcCommunicator) GetControlProgramInfo(id uuid.UUID) ([]opc_custom.ProgramPositionDataType, error) {
	connInfo, err := oc.connector.GetConnectionInfoByUUID(id)
	if err != nil {
		return nil, fmt.Errorf("connection not found: %w", err)
	}
	nodeID := ua.NewNodeIDNumeric(1, 100006)

	machine := interfaces.MachineDataFactory(connInfo.Manufacturer, connInfo.Model)
	if machine == nil {
		return nil, fmt.Errorf("unsupported machine type: %s %s", connInfo.Manufacturer, connInfo.Model)
	}

	val, err := oc.readNodeValue(connInfo.Ctx, connInfo.Conn, nodeID)
	if err != nil {
		log.Printf("Failed to read node %s: %v", nodeID, err)
	}

	// Декодируем значение в структуру
	if err := machine.ConvertNodeToMachineData(nodeID.String(), val); err != nil {
		log.Printf("Failed to convert node %s: %v", nodeID, err)
	}
	data, err := machine.GetExecutionStack()
	if err != nil {
		log.Printf("Failed to get execution stack: %v", err)
	}

	return data, nil
}

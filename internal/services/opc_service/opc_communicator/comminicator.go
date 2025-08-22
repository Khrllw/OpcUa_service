package opc_communicator

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"log"
	"opc_ua_service/internal/interfaces"
)

// OpcCommunicator содержит коннектор с пулом соединений
type OpcCommunicator struct {
	connector interfaces.OpcConnectorService
}

// NewOpcCommunicator создает новый экземпляр OpcCommunicator
func NewOpcCommunicator(connector interfaces.OpcConnectorService) interfaces.OpcCommunicatorService {
	return &OpcCommunicator{
		connector: connector,
	}
}

// CallOPCMethod вызывает метод на OPC UA сервере
func (oc *OpcCommunicator) CallOPCMethod(ctx context.Context, c *client.Client, objectNodeID, methodNodeID ua.NodeID, inputArgs ...ua.Variant) ([]ua.Variant, error) {
	req := &ua.CallRequest{
		MethodsToCall: []ua.CallMethodRequest{
			{
				ObjectID:       objectNodeID,
				MethodID:       methodNodeID,
				InputArguments: inputArgs,
			},
		},
	}

	resp, err := c.Call(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("call request failed: %w", err)
	}

	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("no results returned")
	}

	result := resp.Results[0]
	if !ua.StatusCode.IsGood(result.StatusCode) {
		return nil, fmt.Errorf("method call failed: %s", result.StatusCode)
	}

	return result.OutputArguments, nil
}

// ReadMachineNodes считывает узлы для конкретной машины по sessionID
func (oc *OpcCommunicator) ReadMachineNodes(sessionID string, machineType string) (map[string]*ua.Variant, error) {
	client, err := oc.connector.GetConnectionBySessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection for session %s: %w", sessionID, err)
	}

	// Определяем узлы для чтения в зависимости от типа машины
	var nodes map[string]ua.NodeID
	switch machineType {
	case "CNC_TypeA":
		nodes = map[string]ua.NodeID{
			"SpindleSpeed": ua.NewNodeIDNumeric(3, 1001),
			"FeedRate":     ua.NewNodeIDNumeric(3, 1002),
		}
	case "CNC_TypeB":
		nodes = map[string]ua.NodeID{
			"Temperature": ua.NewNodeIDNumeric(3, 2001),
			"Pressure":    ua.NewNodeIDNumeric(3, 2002),
		}
	default:
		return nil, fmt.Errorf("unknown machine type: %s", machineType)
	}

	results := make(map[string]*ua.Variant)
	for name, nodeID := range nodes {
		res, err := client.Read(nil, &ua.ReadRequest{
			NodesToRead: []ua.ReadValueID{
				{NodeID: nodeID, AttributeID: ua.AttributeIDValue},
			},
		})
		if err != nil || len(res.Results) == 0 {
			results[name] = nil
			continue
		}
		results[name] = &res.Results[0].Value
	}

	return results, nil
}
func (oc *OpcCommunicator) ReadMachineData(sessionID string) (interfaces.MachineData, error) {
	connInfo, err := oc.connector.GetConnectionInfoBySessionID(sessionID)
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

// readNodeValue читает один узел с сервера OPC UA
func (oc *OpcCommunicator) readNodeValue(ctx context.Context, c *client.Client, nodeID ua.NodeIDNumeric) (ua.Variant, error) {
	req := &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{NodeID: nodeID, AttributeID: ua.AttributeIDValue},
		},
	}
	resp, err := c.Read(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("read request failed: %w", err)
	}
	if len(resp.Results) == 0 {
		return nil, fmt.Errorf("no results for node %s", nodeID)
	}
	return resp.Results[0].Value, nil
}

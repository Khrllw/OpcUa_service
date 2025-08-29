package opc_communicator

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
)

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

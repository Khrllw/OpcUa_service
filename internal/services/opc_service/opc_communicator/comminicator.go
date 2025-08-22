package opc_communicator

import (
	"opc_ua_service/internal/interfaces"
)

type OpcCommunicator struct {
}

func NewOpcCommunicator() interfaces.OpcCommunicatorService {
	return &OpcCommunicator{}
}

/* CallOPCMethod вызывает метод на OPC UA сервере
func CallOPCMethod(ctx context.Context, c *client.Client, objectNodeID, methodNodeID ua.NodeID, inputArgs ...ua.Variant) ([]ua.Variant, error) {
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

*/

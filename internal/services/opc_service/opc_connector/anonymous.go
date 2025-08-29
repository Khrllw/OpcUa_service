package opc_connector

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/client"
	connectiion_models "opc_ua_service/internal/domain/models/connection_models"
)

// ConnectAnonymous Анонимное подключение
func (oc *OpcConnector) ConnectAnonymous(config connectiion_models.AnonymousConnection) (*client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	//TODO: Дописать параметры подключения
	clientOpts := []client.Option{
		//client.WithAnonymous(),
	}

	conn, err := oc.createConnection(ctx, config.EndpointURL, clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("ConnectAnonymous failed: %w", err)
	}

	return conn, nil
}

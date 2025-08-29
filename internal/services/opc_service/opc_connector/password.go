package opc_connector

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/client"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
)

// ConnectWithPassword Подключение с логином и паролем
func (oc *OpcConnector) ConnectWithPassword(config connection_models.PasswordConnection) (*client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	//TODO: Дописать параметры подключения
	clientOpts := []client.Option{
		client.WithUserNameIdentity(config.Username, config.Password),
		// Можно добавить SecurityPolicy и Mode через настройки
	}

	conn, err := oc.createConnection(ctx, config.EndpointURL, clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("ConnectWithPassword failed: %w", err)
	}

	return conn, nil
}

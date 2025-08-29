package opc_connector

import (
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/pkg/errors"
)

// GetConnectionByUUID возвращает клиент по UUID
func (oc *OpcConnector) GetConnectionByUUID(id uuid.UUID) (*client.Client, error) {
	oc.mu.RLock()
	info, exists := oc.connections[id]
	oc.mu.RUnlock()
	if !exists {
		return nil, errors.NewNotFoundError("connection not found")
	}

	info.Mu.RLock()
	defer info.Mu.RUnlock()
	if !info.IsHealthy {
		return nil, fmt.Errorf("connection with UUID %s is unhealthy", id)
	}

	return info.Conn, nil
}

// GetConnectionInfoByUUID возвращает ConnectionInfo по UUID
func (oc *OpcConnector) GetConnectionInfoByUUID(id uuid.UUID) (*models.ConnectionInfo, error) {
	oc.mu.RLock()
	info, exists := oc.connections[id]
	oc.mu.RUnlock()
	if !exists {
		return nil, errors.NewNotFoundError("connection not found")
	}

	info.Mu.RLock()
	defer info.Mu.RUnlock()
	if !info.IsHealthy {
		return nil, fmt.Errorf("connection with UUID %s is unhealthy", id)
	}
	return info, nil
}

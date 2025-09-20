package opc_connector

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"time"
)

// Проверка здоровья всех соединений
func (oc *OpcConnector) healthCheckWorker() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	cleanupTicker := time.NewTicker(5 * time.Minute)
	defer cleanupTicker.Stop()

	for {
		select {
		case <-ticker.C:
			oc.checkAllConnectionsHealth()
		case <-cleanupTicker.C:
			oc.Cleanup(10 * time.Minute)
		case <-oc.shutdown:
			return
		}
	}
}

func (oc *OpcConnector) checkAllConnectionsHealth() {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	for uuid, info := range oc.connections {
		info.Mu.Lock()
		if info.IsPolled {
			info.IsHealthy = oc.checkConnectionHealth(info.Conn)
			if !info.IsHealthy {
				info.Mu.Unlock()
				oc.logger.Warn("Connection is unhealthy. Attempting to reconnect...")
				// Попытка переподключиться
				newUUID, err := oc.Reconnect(uuid, info)
				if err != nil || newUUID == nil {
					oc.logger.Warn("Failed to reconnect", "UUID", uuid, "error", err)
				} else {
					oc.logger.Info("Reconnection successful", "UUID", newUUID)
				}
				continue
			}
		}
		info.Mu.Unlock()
	}
}

func (oc *OpcConnector) checkConnectionHealth(conn *client.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := conn.Read(ctx, &ua.ReadRequest{
		NodesToRead: []ua.ReadValueID{
			{NodeID: ua.ObjectIDRootFolder, AttributeID: ua.AttributeIDNodeID},
		},
	})
	return err == nil
}

// recreateConnection создаёт новое подключение по информации из info
func (oc *OpcConnector) Reconnect(ID uuid.UUID, info *models.ConnectionInfo) (*uuid.UUID, error) {

	// Закрываем старое соединение, если есть
	if info.Conn != nil {
		err := oc.CloseConnection(ID)
		if err != nil {
			return nil, err
		}
	}
	info.Mu.Lock()
	defer info.Mu.Unlock()
	if info.Cancel != nil {
		info.Cancel()
	}

	var UUID *uuid.UUID
	var err error

	switch cfg := info.Config.Config.(type) {
	case *connection_models.CertificateConnection:
		UUID, err = oc.CreateCertificateConnection(*cfg)
		if err != nil {
			return nil, err
		}
	case *connection_models.PasswordConnection:
		UUID, err = oc.CreatePasswordConnection(*cfg)
		if err != nil {
			return nil, err
		}
	case *connection_models.AnonymousConnection:
		UUID, err = oc.CreateAnonymousConnection(*cfg)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("undefined connection type: %T", info.Config)
	}
	if UUID != nil {
		return UUID, nil
	}
	return nil, fmt.Errorf("empty connection UUID")
}

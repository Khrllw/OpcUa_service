package opc_connector

import (
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/models"
	"sync/atomic"
)

// GetAllConnectionsInfo возвращает информацию о всех соединениях
func (oc *OpcConnector) GetAllConnectionsInfo() map[uuid.UUID]*models.ConnectionInfo {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	all := make([]*models.ConnectionInfo, 0, len(oc.connections))
	for _, info := range oc.connections {
		all = append(all, info)
	}
	return oc.connections
}

// GetGlobalStats возвращает краткую статистику по соединениям
func (oc *OpcConnector) GetGlobalStats() models.ConnectorStats {
	return models.ConnectorStats{
		TotalConnections:  atomic.LoadInt64(&oc.stats.TotalConnections),
		FailedConnections: atomic.LoadInt64(&oc.stats.FailedConnections),
		ActiveConnections: atomic.LoadInt64(&oc.stats.ActiveConnections),
		PoolSize:          atomic.LoadInt64(&oc.stats.PoolSize),
	}
}

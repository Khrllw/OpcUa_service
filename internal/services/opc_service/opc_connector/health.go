package opc_connector

import (
	"context"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"time"
)

// Проверка здоровья всех соединений
func (oc *OpcConnector) healthCheckWorker() {
	ticker := time.NewTicker(1 * time.Minute)
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

	for _, info := range oc.connections {
		info.Mu.Lock()
		info.IsHealthy = oc.checkConnectionHealth(info.Conn)
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

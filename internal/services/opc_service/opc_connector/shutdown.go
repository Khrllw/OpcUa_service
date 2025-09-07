package opc_connector

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/models"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// SetupSignalHandler ловит SIGINT и SIGTERM
func (oc *OpcConnector) SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigs:
			oc.logger.Info("Signal %s received, shutting down...", sig)
			cancel()
			oc.Shutdown()
		case <-oc.shutdown:
			cancel()
		}
	}()

	return ctx
}

func (oc *OpcConnector) Shutdown() {
	close(oc.shutdown)
	oc.CloseAll()
}

// CloseAll закрывает все подключения
func (oc *OpcConnector) CloseAll() {
	// Копируем ссылки на все соединения под Lock
	oc.mu.Lock()
	conns := make([]*models.ConnectionInfo, 0, len(oc.connections))
	for _, info := range oc.connections {
		conns = append(conns, info)
	}
	oc.connections = make(map[uuid.UUID]*models.ConnectionInfo)
	atomic.StoreInt64(&oc.stats.ActiveConnections, 0)
	atomic.StoreInt64(&oc.stats.PoolSize, 0)
	oc.mu.Unlock()

	// Закрываем все соединения вне Lock
	for _, info := range conns {
		info.Mu.Lock()
		ctxClose, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		if err := info.Conn.Close(ctxClose); err != nil {
			oc.logger.Error("Failed to close connection %s: %v", info.SessionID, err)
		}
		cancel()
		info.Cancel()
		info.Mu.Unlock()
	}
}

// CloseConnection закрывает подключение и удаляет его из пула
func (oc *OpcConnector) CloseConnection(id uuid.UUID) error {
	oc.mu.Lock()
	info, exists := oc.connections[id]
	if !exists || info == nil {
		oc.mu.Unlock()
		return fmt.Errorf("connection with UUID %s not found", id)
	}
	delete(oc.connections, id)
	atomic.AddInt64(&oc.stats.ActiveConnections, -1)
	atomic.AddInt64(&oc.stats.PoolSize, -1)
	oc.mu.Unlock()

	// Закрываем соединение вне Lock
	info.Mu.Lock()
	defer info.Mu.Unlock()

	ctxClose, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := info.Conn.Close(ctxClose); err != nil {
		return fmt.Errorf("failed to close connection %s: %v", info.SessionID, err)
	}

	info.Cancel()

	oc.logger.Info("Connection with UUID %s closed successfully", id)
	return nil
}

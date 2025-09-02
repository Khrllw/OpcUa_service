package opc_connector

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"github.com/google/uuid"
	"log"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	_ "opc_ua_service/internal/domain/models/opc_custom"
	"opc_ua_service/internal/interfaces"
	"sync"
	"sync/atomic"
	"time"
)

type OpcConnector struct {
	mu          sync.RWMutex
	connections map[uuid.UUID]*models.ConnectionInfo // UUID -> connection info
	stats       models.ConnectorStats
	shutdown    chan struct{}
	certManager interfaces.CertificateManagerService
}

func NewOpcConnector(certManager interfaces.CertificateManagerService) interfaces.OpcConnectorService {
	connector := &OpcConnector{
		connections: make(map[uuid.UUID]*models.ConnectionInfo),
		shutdown:    make(chan struct{}),
		certManager: certManager,
	}
	go connector.healthCheckWorker()
	return connector
}

// GenerateUUID возвращает новую UUID в виде строки
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// -------------------------------------- PUBLIC --------------------------------------

// FindOpenConnection ищет уже существующее соединение по UUID
func (oc *OpcConnector) FindOpenConnection(id uuid.UUID) *models.ConnectionInfo {
	oc.mu.Lock()
	info, exists := oc.connections[id]
	if !exists {
		oc.mu.Unlock()
		return nil
	}
	oc.mu.Unlock()

	info.Mu.Lock()
	defer info.Mu.Unlock()

	if !oc.checkConnectionHealth(info.Conn) {
		oc.closeConnectionInternal(id)

		oc.mu.Lock()
		delete(oc.connections, id)
		oc.mu.Unlock()

		atomic.AddInt64(&oc.stats.ActiveConnections, -1)
		atomic.AddInt64(&oc.stats.PoolSize, -1)
		return nil
	}

	info.LastUsed = time.Now()
	atomic.AddInt64(&info.UseCount, 1)

	return info
}

func (oc *OpcConnector) CreateAnonymousConnection(config connection_models.AnonymousConnection) (*uuid.UUID, error) {
	return nil, nil
}

func (oc *OpcConnector) CreatePasswordConnection(config connection_models.PasswordConnection) (*uuid.UUID, error) {
	return nil, nil
}

// CreateConnection создаёт новое подключение и сохраняет его по UUID
func (oc *OpcConnector) CreateCertificateConnection(config connection_models.CertificateConnection) (*uuid.UUID, error) {
	// Создаем новое подключение
	conn, err := oc.ConnectWithCertificate(config)
	if err != nil {
		atomic.AddInt64(&oc.stats.FailedConnections, 1)
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sessionID := conn.SessionID().(ua.NodeIDNumeric).String()

	cfg := connection_models.ConnectionConfig{
		Config: &config,
	}

	info := &models.ConnectionInfo{
		Conn:         conn,
		Ctx:          ctx,
		Cancel:       cancel,
		SessionID:    sessionID,
		Config:       cfg,
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
		IsHealthy:    true,
		UseCount:     1,
		Mu:           sync.RWMutex{},
		Manufacturer: cfg.Config.GetManufacturer(), // <-- берем из конфига
		Model:        cfg.Config.GetModel(),
	}

	oc.mu.Lock()
	defer oc.mu.Unlock()

	id := GenerateUUID()
	if _, exists := oc.connections[id]; exists {
		// Если по какой-то причине соединение уже есть, закрываем новое
		info.Cancel()
		_ = conn.Close(ctx)
		return &id, nil
	}

	oc.connections[id] = info
	atomic.AddInt64(&oc.stats.TotalConnections, 1)
	atomic.AddInt64(&oc.stats.PoolSize, 1)
	atomic.AddInt64(&oc.stats.ActiveConnections, 1)

	return &id, nil
}

// GetConnection получает подключение по конфигу
func (oc *OpcConnector) GetConnection(config connection_models.CertificateConnection) (*client.Client, error) {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	for _, info := range oc.connections {
		info.Mu.Lock()
		if info.IsHealthy {
			info.LastUsed = time.Now()
			info.UseCount++
			conn := info.Conn
			info.Mu.Unlock()
			return conn, nil
		}
		info.Mu.Unlock()
	}

	return nil, fmt.Errorf("connection not found in pool")
}

// Cleanup удаляет неиспользуемые соединения
func (oc *OpcConnector) Cleanup(maxIdleTime time.Duration) int {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	count := 0
	now := time.Now()
	for id, info := range oc.connections {
		info.Mu.RLock()
		lastUsed := info.LastUsed
		useCount := info.UseCount
		info.Mu.RUnlock()

		if now.Sub(lastUsed) > maxIdleTime && useCount == 0 {
			err := oc.closeConnectionInternal(id)
			if err != nil {
				return 0
			}
			delete(oc.connections, id)
			count++
			atomic.AddInt64(&oc.stats.ActiveConnections, -1)
			atomic.AddInt64(&oc.stats.PoolSize, -1)
		}
	}
	return count
}

// ---------------------------------------------------------------------------
// closeConnectionInternal Функция закрытия соединения
func (oc *OpcConnector) closeConnectionInternal(id uuid.UUID) error {
	// Получаем соединение из пула
	oc.mu.RLock()
	info, exists := oc.connections[id]
	oc.mu.RUnlock()

	if !exists || info == nil {
		return fmt.Errorf("connection not found")
	}

	info.Mu.Lock()
	defer info.Mu.Unlock()

	ctxClose, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := info.Conn.Close(ctxClose); err != nil {
		return fmt.Errorf("failed to close connection %s: %v", info.SessionID, err)
	}

	// Отменяем контекст
	info.Cancel()

	// Удаляем из пула
	oc.mu.Lock()
	delete(oc.connections, id)
	oc.mu.Unlock()
	return nil
}

// createConnection Общая функция подключения
func (oc *OpcConnector) createConnection(ctx context.Context, endpoint string, opts ...client.Option) (*client.Client, error) {
	conn, err := client.Dial(ctx, endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	log.Printf("Successfully connected to OPC UA server: %s", endpoint)
	return conn, nil
}

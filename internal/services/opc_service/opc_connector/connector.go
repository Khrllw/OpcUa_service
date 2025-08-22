package opc_connector

import (
	"context"
	"fmt"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"log"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/internal/interfaces"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

type OpcConnector struct {
	mu          sync.RWMutex
	connections map[string]*models.ConnectionInfo // sessionID -> connection info
	stats       models.ConnectorStats
	shutdown    chan struct{}
	certManager interfaces.CertificateManagerService
}

func NewOpcConnector(certManager interfaces.CertificateManagerService) interfaces.OpcConnectorService {
	connector := &OpcConnector{
		connections: make(map[string]*models.ConnectionInfo),
		shutdown:    make(chan struct{}),
		certManager: certManager,
	}
	go connector.healthCheckWorker()
	return connector
}

// ------------------- Public Methods -------------------

// FindOpenConnectionByConn ищет уже существующее соединение по объекту client.Client
func (oc *OpcConnector) FindOpenConnectionByConn(conn *client.Client) *models.ConnectionInfo {
	if conn == nil {
		return nil
	}
	sessionID := conn.SessionID().(ua.NodeIDNumeric).String()

	oc.mu.Lock()
	defer oc.mu.Unlock()

	info, exists := oc.connections[sessionID]
	if !exists {
		return nil
	}

	info.Mu.Lock()
	defer info.Mu.Unlock()

	if !oc.checkConnectionHealth(info.Conn) {
		oc.closeConnectionInternal(info)
		delete(oc.connections, sessionID)
		atomic.AddInt64(&oc.stats.ActiveConnections, -1)
		atomic.AddInt64(&oc.stats.PoolSize, -1)
		return nil
	}

	info.LastUsed = time.Now()
	atomic.AddInt64(&info.UseCount, 1)
	return info
}

// CreateConnection создаёт новое подключение и сохраняет его по sessionID
func (oc *OpcConnector) CreateConnection(config models.ConnectionConfig) (*models.ConnectionInfo, error) {
	// Сначала ищем уже существующее подключение с тем же конфигом
	for _, info := range oc.connections {
		info.Mu.RLock()
		sameConfig := info.Config == config && info.IsHealthy
		info.Mu.RUnlock()
		if sameConfig {
			info.Mu.Lock()
			info.UseCount++
			info.LastUsed = time.Now()
			info.Mu.Unlock()
			return info, nil
		}
	}

	// Создаем новое подключение
	conn, err := oc.createConnection(config)
	if err != nil {
		atomic.AddInt64(&oc.stats.FailedConnections, 1)
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	sessionID := conn.SessionID().(ua.NodeIDNumeric).String()

	info := &models.ConnectionInfo{
		Conn:         conn,
		Ctx:          ctx,
		Cancel:       cancel,
		SessionID:    sessionID,
		Config:       config,
		CreatedAt:    time.Now(),
		LastUsed:     time.Now(),
		IsHealthy:    true,
		UseCount:     1,
		Mu:           sync.RWMutex{},
		Manufacturer: config.Manufacturer, // <-- берем из конфига
		Model:        config.Model,
	}

	oc.mu.Lock()
	defer oc.mu.Unlock()
	if _, exists := oc.connections[sessionID]; exists {
		// Если по какой-то причине соединение уже есть, закрываем новое
		info.Cancel()
		_ = conn.Close(ctx)
		return oc.connections[sessionID], nil
	}

	oc.connections[sessionID] = info
	atomic.AddInt64(&oc.stats.TotalConnections, 1)
	atomic.AddInt64(&oc.stats.PoolSize, 1)
	atomic.AddInt64(&oc.stats.ActiveConnections, 1)

	return info, nil
}

// GetOrCreateConnection получает или создаёт подключение по конфигу
func (oc *OpcConnector) GetOrCreateConnection(config models.ConnectionConfig) (*client.Client, error) {
	oc.mu.RLock()
	for _, info := range oc.connections {
		info.Mu.RLock()
		if info.Config == config && info.IsHealthy {
			info.Mu.RUnlock()
			oc.mu.RUnlock()

			info.Mu.Lock()
			info.LastUsed = time.Now()
			info.UseCount++
			info.Mu.Unlock()
			return info.Conn, nil
		}
		info.Mu.RUnlock()
	}
	oc.mu.RUnlock()

	info, err := oc.CreateConnection(config)
	if err != nil {
		return nil, err
	}
	return info.Conn, nil
}

// CloseConnectionByConn закрывает соединение по объекту client
func (oc *OpcConnector) CloseConnectionByConn(conn *client.Client) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}
	sessionID := conn.SessionID().(ua.NodeIDNumeric).String()

	oc.mu.Lock()
	defer oc.mu.Unlock()

	info, exists := oc.connections[sessionID]
	if !exists {
		return fmt.Errorf("connection not found")
	}

	oc.closeConnectionInternal(info)
	delete(oc.connections, sessionID)
	atomic.AddInt64(&oc.stats.ActiveConnections, -1)
	atomic.AddInt64(&oc.stats.PoolSize, -1)
	return nil
}

// CloseAll закрывает все подключения
func (oc *OpcConnector) CloseAll() {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	for _, info := range oc.connections {
		oc.closeConnectionInternal(info)
	}
	oc.connections = make(map[string]*models.ConnectionInfo)
	atomic.StoreInt64(&oc.stats.ActiveConnections, 0)
	atomic.StoreInt64(&oc.stats.PoolSize, 0)
}

// GetConnectionBySessionID возвращает клиент по sessionID
func (oc *OpcConnector) GetConnectionBySessionID(sessionID string) (*client.Client, error) {
	oc.mu.RLock()
	info, exists := oc.connections[sessionID]
	oc.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("connection not found")
	}

	info.Mu.RLock()
	defer info.Mu.RUnlock()
	if !info.IsHealthy {
		return nil, fmt.Errorf("connection with sessionID %s is unhealthy", sessionID)
	}

	return info.Conn, nil
}

// GetConnectionInfoBySessionID возвращает ConnectionInfo по sessionID
func (oc *OpcConnector) GetConnectionInfoBySessionID(sessionID string) (*models.ConnectionInfo, error) {
	oc.mu.RLock()
	info, exists := oc.connections[sessionID]
	oc.mu.RUnlock()
	if !exists {
		return nil, fmt.Errorf("connection not found")
	}

	info.Mu.RLock()
	defer info.Mu.RUnlock()
	if !info.IsHealthy {
		return nil, fmt.Errorf("connection with sessionID %s is unhealthy", sessionID)
	}
	return info, nil
}

// GetAllConnectionsInfo возвращает все соединения
func (oc *OpcConnector) GetAllConnectionsInfo() []*models.ConnectionInfo {
	oc.mu.RLock()
	defer oc.mu.RUnlock()

	all := make([]*models.ConnectionInfo, 0, len(oc.connections))
	for _, info := range oc.connections {
		all = append(all, info)
	}
	return all
}

// Cleanup удаляет неиспользуемые соединения
func (oc *OpcConnector) Cleanup(maxIdleTime time.Duration) int {
	oc.mu.Lock()
	defer oc.mu.Unlock()

	count := 0
	now := time.Now()
	for sessionID, info := range oc.connections {
		info.Mu.RLock()
		lastUsed := info.LastUsed
		useCount := info.UseCount
		info.Mu.RUnlock()

		if now.Sub(lastUsed) > maxIdleTime && useCount == 0 {
			oc.closeConnectionInternal(info)
			delete(oc.connections, sessionID)
			count++
			atomic.AddInt64(&oc.stats.ActiveConnections, -1)
			atomic.AddInt64(&oc.stats.PoolSize, -1)
		}
	}
	return count
}

// ------------------- Private Methods -------------------

func (oc *OpcConnector) createConnection(config models.ConnectionConfig) (*client.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	certBytes, clientCert, clientKey := oc.certManager.LoadClientCredentials(config.Certificate, config.Key)
	if err := oc.certManager.VerifyKeyMatchesCert(clientCert, clientKey); err != nil {
		return nil, fmt.Errorf("certificate verification failed: %w", err)
	}

	endpoint, policyID := oc.certManager.SelectCertificateEndpoint(ctx, config.EndpointURL)
	clientOpts := oc.certManager.BuildClientOptions(endpoint, policyID, certBytes, clientKey)

	conn, err := client.Dial(ctx, config.EndpointURL, clientOpts...)
	if err != nil {
		return nil, fmt.Errorf("dial failed: %w", err)
	}

	log.Printf("Successfully connected to OPC UA server: %s", config.EndpointURL)
	return conn, nil
}

func (oc *OpcConnector) closeConnectionInternal(info *models.ConnectionInfo) {
	ctxClose, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := info.Conn.Close(ctxClose); err != nil {
		log.Printf("Failed to close connection: %v", err)
	} else {
		log.Printf("Connection closed: %s", info.Config.EndpointURL)
	}
	info.Cancel()
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

func (oc *OpcConnector) SetupSignalHandler() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		select {
		case sig := <-sigs:
			log.Printf("Signal %s received, shutting down...", sig)
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

// CloseConnectionBySessionID закрывает соединение по sessionID и удаляет его из пула
func (oc *OpcConnector) CloseConnectionBySessionID(sessionID string) error {
	if sessionID == "" {
		return fmt.Errorf("sessionID is empty")
	}

	oc.mu.Lock()
	defer oc.mu.Unlock()

	info, exists := oc.connections[sessionID]
	if !exists {
		return fmt.Errorf("connection with sessionID %s not found", sessionID)
	}

	// Закрываем соединение безопасно
	oc.closeConnectionInternal(info)
	delete(oc.connections, sessionID)

	// Обновляем статистику
	atomic.AddInt64(&oc.stats.ActiveConnections, -1)
	atomic.AddInt64(&oc.stats.PoolSize, -1)

	log.Printf("Connection with sessionID %s closed successfully", sessionID)
	return nil
}

// CloseConnection закрывает подключение и удаляет его из пула
func (oc *OpcConnector) CloseConnection(conn *client.Client) error {
	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	sessionID := conn.SessionID().(ua.NodeIDNumeric).String()

	oc.mu.Lock()
	defer oc.mu.Unlock()

	info, exists := oc.connections[sessionID]
	if !exists {
		return fmt.Errorf("connection with sessionID %s not found", sessionID)
	}

	// Закрываем соединение безопасно
	oc.closeConnectionInternal(info)
	delete(oc.connections, sessionID)

	// Обновляем статистику
	atomic.AddInt64(&oc.stats.ActiveConnections, -1)
	atomic.AddInt64(&oc.stats.PoolSize, -1)

	log.Printf("Connection with sessionID %s closed successfully", sessionID)
	return nil
}

// GetGlobalStats возвращает глобальную статистику по всем соединениям
func (oc *OpcConnector) GetGlobalStats() models.ConnectorStats {
	return models.ConnectorStats{
		TotalConnections:  atomic.LoadInt64(&oc.stats.TotalConnections),
		FailedConnections: atomic.LoadInt64(&oc.stats.FailedConnections),
		ActiveConnections: atomic.LoadInt64(&oc.stats.ActiveConnections),
		PoolSize:          atomic.LoadInt64(&oc.stats.PoolSize),
	}
}

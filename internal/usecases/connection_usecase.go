package usecases

import (
	"errors"
	"fmt"
	"github.com/awcullen/opcua/ua"
	"log"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/internal/interfaces"
	"time"
)

type ConnectionUsecase struct {
	OpcService interfaces.OpcService
}

func (u *ConnectionUsecase) GetConnectionStats(sessionID string) (map[string]interface{}, error) {
	// Получаем соединение по sessionID через сервис OPC
	connInfo, err := u.OpcService.GetConnectionInfoBySessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find connection: %w", err)
	}

	// Блокируем данные соединения для безопасного чтения
	connInfo.Mu.RLock()
	defer connInfo.Mu.RUnlock()

	stats := map[string]interface{}{
		"session_id": connInfo.SessionID,
		"endpoint":   connInfo.Config.EndpointURL,
		"created_at": connInfo.CreatedAt,
		"last_used":  connInfo.LastUsed,
		"use_count":  connInfo.UseCount,
		"is_healthy": connInfo.IsHealthy,
	}

	return stats, nil
}

func NewConnectionUsecase(s interfaces.OpcService) *ConnectionUsecase {
	return &ConnectionUsecase{s}
}

func (u *ConnectionUsecase) LoginClientAnonymous() (string, error) {
	// Генерация уникального ID и токена для анонимного клиента
	return "jwt-token-abc", nil
}

func (u *ConnectionUsecase) LoginClientPassword(request models.ConnectionRequest) (string, error) {
	// Проверка username/password в базе
	if request.Username == "client1" && request.Password == "secret" {
		return "jwt-token-xyz", nil
	}
	return "", errors.New("invalid credentials")
}

func (u *ConnectionUsecase) LoginClientCertificate(request models.ConnectionRequest) (string, error) {
	/* Контекст с обработкой сигналов
	ctx := u.OpcConnector.SetupSignalHandler()

	// Загрузка сертификатов
	certBytes, clientCert, clientKey := u.CertManager.LoadClientCredentials(request.Certificate, request.Key)
	serverCert := u.CertManager.LoadServerCertificate("certs/new_server/heopcua-rootca_khrllw--340595_2025aug06_132042.der")

	// Проверка соответствия ключа сертификату
	if err := u.CertManager.VerifyKeyMatchesCert(clientCert, clientKey); err != nil {
		return "", fmt.Errorf("key does not match certificate: %w", err)
	}

	u.CertManager.PrintCertInfo("Client", clientCert)
	u.CertManager.PrintCertInfo("Server", serverCert)

	// Выбор endpoint и политики (пример выбора первого)
	endpoint, policyID := u.CertManager.SelectCertificateEndpoint(ctx, "opc.tcp://KHRLLW_-340595:4840/HEIDENHAIN/NC")

	// Настройка клиентских опций
	clientOpts := u.CertManager.BuildClientOptions(endpoint, policyID, certBytes, clientKey)

	// Подключение к OPC UA серверу
	conn := u.OpcConnector.Connect(ctx, endpoint.EndpointURL, clientOpts)
	if conn == nil {
		return "", errors.New("failed to connect to OPC UA server")
	}
	defer u.OpcConnector.CloseConnection(conn)

	sessionID := conn.SessionID().(ua.NodeIDNumeric)
	return sessionID.String(), nil
	*/
	return "", nil
}

func (u *ConnectionUsecase) ConnectByCert(request models.ConnectionRequest) (*models.ConnectionInfo, error) {
	// Проверяем обязательные поля
	if request.EndpointURL == "" {
		return nil, fmt.Errorf("endpointURL is required")
	}
	if request.Certificate == "" {
		return nil, fmt.Errorf("certificate is required")
	}
	if request.Key == "" {
		return nil, fmt.Errorf("key is required")
	}

	// Настраиваем значения по умолчанию
	policy := request.Policy
	if policy == "" {
		policy = models.PolicyBasic256Sha256
	}
	mode := request.Mode
	if mode == "" {
		mode = models.ModeSignAndEncrypt
	}

	timeout := request.Timeout
	if timeout == 0 {
		timeout = 30
	}

	// Проверяем корректность enum
	if err := models.SecurityPolicyEnum(policy).Validate(); err != nil {
		return nil, fmt.Errorf("invalid security policy: %w", err)
	}
	if err := models.MessageSecurityModeEnum(mode).Validate(); err != nil {
		return nil, fmt.Errorf("invalid message security mode: %w", err)
	}

	// Создаем конфиг
	config := models.ConnectionConfig{
		EndpointURL: request.EndpointURL,
		Certificate: request.Certificate,
		Key:         request.Key,
		Policy:      string(policy),
		Mode:        string(mode),
		Timeout:     time.Duration(timeout) * time.Second,
	}

	// Получаем или создаем подключение через пул
	conn, err := u.OpcService.CreateConnection(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	sessionID := conn.Conn.SessionID().(ua.NodeIDNumeric)
	log.Printf("Connected successfully. Session ID: %s", sessionID)

	return conn, nil
}

// DisconnectBySessionID закрывает соединение по session ID
func (u *ConnectionUsecase) DisconnectBySessionID(sessionID string) error {
	conn, err := u.OpcService.GetConnectionBySessionID(sessionID)
	if err != nil {
		return fmt.Errorf("failed to find connection: %w", err)
	}

	if err := u.OpcService.CloseConnection(conn); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	log.Printf("Successfully closed connection with session ID: %s", sessionID)
	return nil
}

// DisconnectAll закрывает все соединения
func (u *ConnectionUsecase) DisconnectAll() (int, error) {
	stats := u.OpcService.GetGlobalStats()
	activeConnections := stats.ActiveConnections

	u.OpcService.CloseAll()

	log.Printf("Closed all %d active connections", activeConnections)
	return int(activeConnections), nil
}

// GetActiveConnections возвращает список активных соединений
func (u *ConnectionUsecase) GetActiveConnections() []*models.ConnectionInfo {
	// Получаем все соединения из сервиса
	connectionsInfo := u.OpcService.GetAllConnectionsInfo() // предполагаем, что метод возвращает map[string]*ConnectionInfo

	return connectionsInfo
}

// CleanupIdleConnections очищает неиспользуемые соединения
func (u *ConnectionUsecase) CleanupIdleConnections(maxIdleMinutes int) int {
	cleaned := u.OpcService.Cleanup(time.Duration(maxIdleMinutes) * time.Minute)
	log.Printf("Cleaned up %d idle connections (idle time > %d minutes)", cleaned, maxIdleMinutes)
	return cleaned
}

/*
// ConnectAndMonitor устанавливает подключение и начинает мониторинг
func (u *ConnectionUsecase) ConnectAndMonitor(request models.ConnectionRequest) (string, error) {
	config := models.ConnectionConfig{
		EndpointURL: request.EndpointURL,
		Certificate: request.Certificate,
		Key:         request.Key,
		Policy:      "Basic256Sha256",
		Mode:        "SignAndEncrypt",
		Timeout:     30 * time.Second,
	}

	// Получаем или создаем подключение через пул
	conn, err := u.opcConnector.GetOrCreateConnection(config)
	if err != nil {
		return "", fmt.Errorf("failed to get connection: %w", err)
	}

	// Добавляем узлы для мониторинга
	u.configureMonitoringNodes()

	// Запускаем мониторинг с использованием контекста с обработкой сигналов
	ctx := u.opcConnector.SetupSignalHandler()
	u.opcCommunicator.Start(ctx, conn)

	sessionID := conn.SessionID().String()
	log.Printf("Connected successfully. Session ID: %s", sessionID)

	return sessionID, nil
}

// configureMonitoringNodes настраивает узлы для мониторинга
func (u *ConnectionUsecase) configureMonitoringNodes() {
	// Добавляем узлы для чтения
	nodes := []string{
		"ns=2;s=Machine/Temperature",
		"ns=2;s=Machine/Pressure",
		"ns=2;s=Machine/Status",
		"ns=2;s=Machine/ProductionRate",
		"ns=2;s=Machine/EnergyConsumption",
	}

	for _, node := range nodes {
		u.opcCommunicator.AddReadingNode(node)
	}
	log.Printf("Configured %d nodes for monitoring", len(nodes))
}


*/

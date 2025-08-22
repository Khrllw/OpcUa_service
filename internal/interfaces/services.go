package interfaces

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"opc_ua_service/internal/domain/models"
	"time"
)

type OpcService interface {
	CertificateManagerService
	OpcConnectorService
	OpcCommunicatorService
}

type CertificateManagerService interface {
	LoadCertificate(certPath string) ([]byte, *x509.Certificate, error)
	SelectCertificateEndpoint(ctx context.Context, endpointURL string) (*ua.EndpointDescription, string)
	PrintCertInfo(label string, cert *x509.Certificate)
	VerifyKeyMatchesCert(cert *x509.Certificate, key *rsa.PrivateKey) error
	LoadPrivateKey(keyPath string) (*rsa.PrivateKey, error)
	LoadServerCertificate(certPath string) *x509.Certificate
	LoadClientCredentials(certPath, keyPath string) ([]byte, *x509.Certificate, *rsa.PrivateKey)
	BuildClientOptions(endpoint *ua.EndpointDescription, policyID string, certBytes []byte, key *rsa.PrivateKey) []client.Option
}

type OpcConnectorService interface {
	// Получить или создать подключение по конфигу
	GetOrCreateConnection(config models.ConnectionConfig) (*client.Client, error)

	// Создать новое подключение и вернуть ConnectionInfo
	CreateConnection(config models.ConnectionConfig) (*models.ConnectionInfo, error)

	// Найти открытое и здоровое соединение по объекту client
	FindOpenConnectionByConn(conn *client.Client) *models.ConnectionInfo

	// Закрыть подключение по объекту client
	CloseConnectionByConn(conn *client.Client) error

	// Закрыть все соединения
	CloseAll()

	// Получить клиент по sessionID
	GetConnectionBySessionID(sessionID string) (*client.Client, error)

	// Получить ConnectionInfo по sessionID
	GetConnectionInfoBySessionID(sessionID string) (*models.ConnectionInfo, error)

	// Получить все соединения
	GetAllConnectionsInfo() []*models.ConnectionInfo

	// Очистка неиспользуемых соединений
	Cleanup(maxIdleTime time.Duration) int

	// CloseConnection закрывает подключение и удаляет его из пула
	CloseConnection(conn *client.Client) error

	// GetGlobalStats возвращает глобальную статистику по всем соединениям
	GetGlobalStats() models.ConnectorStats
}

type OpcCommunicatorService interface {
	CallOPCMethod(ctx context.Context, c *client.Client, objectNodeID, methodNodeID ua.NodeID, inputArgs ...ua.Variant) ([]ua.Variant, error)
	ReadMachineNodes(sessionID string, machineType string) (map[string]*ua.Variant, error)
	ReadMachineData(sessionID string) (MachineData, error)
}

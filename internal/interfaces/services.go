package interfaces

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"github.com/awcullen/opcua/client"
	"github.com/awcullen/opcua/ua"
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/internal/domain/models/opc_custom"
	"time"
)

type OpcService interface {
	CertificateManagerService
	OpcConnectorService
	OpcCommunicatorService
}

type CertificateManagerService interface {
	SelectCertificateEndpoint(ctx context.Context, endpointURL string) (*ua.EndpointDescription, string)
	PrintCertInfo(label string, cert *x509.Certificate)
	VerifyKeyMatchesCert(cert *x509.Certificate, key *rsa.PrivateKey) error

	LoadPrivateKey(keyPath string) ([]byte, *rsa.PrivateKey, error)
	LoadPrivateKeyBytes(data []byte) (*rsa.PrivateKey, error)

	LoadCertificate(certPath string) ([]byte, *x509.Certificate, error)
	LoadServerCertificate(certPath string) *x509.Certificate
	LoadCertificateBytes(data []byte) (*x509.Certificate, error)

	LoadClientCredentials(certPath, keyPath string) ([]byte, *x509.Certificate, *rsa.PrivateKey)
	LoadClientCredentialsBytes(certBytes, keyBytes []byte) ([]byte, *x509.Certificate, *rsa.PrivateKey)
	LoadClientCredentialsBase64(certBase64, keyBase64 string) ([]byte, *x509.Certificate, *rsa.PrivateKey)
	BuildClientOptions(endpoint *ua.EndpointDescription, policyID string, certBytes []byte, key *rsa.PrivateKey) []client.Option
}

type OpcConnectorService interface {
	GetConnection(config connection_models.CertificateConnection) (*client.Client, error)
	CreateAnonymousConnection(config connection_models.AnonymousConnection) (*uuid.UUID, error)
	CreatePasswordConnection(config connection_models.PasswordConnection) (*uuid.UUID, error)
	CreateCertificateConnection(config connection_models.CertificateConnection) (*uuid.UUID, error)
	CloseAll()
	GetConnectionByUUID(id uuid.UUID) (*client.Client, error)
	GetConnectionInfoByUUID(id uuid.UUID) (*models.ConnectionInfo, error)
	//GetControlProgramInfo(sessionID string) ([]opc_custom.ProgramPositionDataType, error)

	Cleanup(maxIdleTime time.Duration) int
	CloseConnection(id uuid.UUID) error
	GetGlobalStats() models.ConnectorStats
	GetAllConnectionsInfo() map[uuid.UUID]*models.ConnectionInfo
	FindOpenConnection(id uuid.UUID) *models.ConnectionInfo

	Base64ToBytes(base64Str string) ([]byte, error)
	BytesToBase64(data []byte) string
}

type OpcCommunicatorService interface {
	CallOPCMethod(ctx context.Context, c *client.Client, objectNodeID, methodNodeID ua.NodeID, inputArgs ...ua.Variant) ([]ua.Variant, error)
	ReadMachineData(id uuid.UUID) (MachineData, error)
	GetControlProgramInfo(id uuid.UUID) ([]opc_custom.ProgramPositionDataType, error)
	StartPollingForMachine(id uuid.UUID) error
	StopPollingForMachine(id uuid.UUID) error
}

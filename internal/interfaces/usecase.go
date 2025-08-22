package interfaces

import (
	"opc_ua_service/internal/domain/models"
)

type Usecases interface {
	ConnectionUsecase
	PoolingUsecase
}

type ConnectionUsecase interface {
	LoginClientAnonymous() (token string, err error)
	LoginClientPassword(request models.ConnectionRequest) (token string, err error)
	LoginClientCertificate(request models.ConnectionRequest) (string, error)
	ConnectByCert(request models.ConnectionRequest) (*models.ConnectionInfo, error)
	DisconnectBySessionID(sessionID string) error
	DisconnectAll() (int, error)
	GetConnectionStats(sessionID string) (map[string]interface{}, error)
	GetActiveConnections() []*models.ConnectionInfo
	CleanupIdleConnections(maxIdleMinutes int) int
}

type PoolingUsecase interface {
	StartPooling() error
	StopPooling()
}

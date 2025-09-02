package interfaces

import (
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/pkg/errors"
)

type Usecases interface {
	ConnectionUsecase
	PollingUsecase
}

type ConnectionUsecase interface {
	ConnectAnonymous(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError)
	ConnectWithPassword(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError)
	ConnectWithCertificate(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError)

	NewCertificateConnectionFromRequest(req *models.ConnectionRequest) (connection_models.CertificateConnection, error)

	DisconnectByUUID(id uuid.UUID) (*bool, *errors.AppError)
	DisconnectAll() (int, *errors.AppError)

	GetActiveConnections() models.ConnectionPoolResponse

	GetConnectionState(id uuid.UUID) (*models.ConnectionInfoResponse, *errors.AppError)
	CleanupIdleConnections(maxIdleMinutes int) int
	RestoreConnection(machine entities.CncMachine) (*models.ConnectionInfo, *errors.AppError)
}

type PollingUsecase interface {
	GetControlProgram(req models.GetControlProgramRequest) (*models.ControlProgramInfoRequest, error)

	StartPollingMachine(machineID uuid.UUID) *errors.AppError
	StopPollingMachine(machineID uuid.UUID) *errors.AppError
}

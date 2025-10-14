package interfaces

import (
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/internal/domain/models/connection"
	"opc_ua_service/pkg/errors"
)

type Usecases interface {
	ConnectionUsecase
	PollingUsecase
}

type ConnectionUsecase interface {
	ConnectAnonymous(request connection.ConnectionRequest) (entities.CncMachine, *errors.AppError)
	ConnectWithPassword(request connection.ConnectionRequest) (entities.CncMachine, *errors.AppError)
	ConnectWithCertificate(request connection.ConnectionRequest) (entities.CncMachine, *errors.AppError)
	RestoreConnection(machine entities.CncMachine) (*models.ConnectionInfo, *errors.AppError)

	DisconnectByID(machineID uint) (*bool, *errors.AppError)
	DisconnectAll() (int, *errors.AppError)
	CleanupIdleConnections(maxIdleMinutes int) int

	GetActiveConnections() connection.ConnectionPoolResponse
	GetConnectionState(machineID uint) (*connection.ConnectionInfoResponse, *errors.AppError)
}

type PollingUsecase interface {
	GetControlProgram(req connection.GetControlProgramRequest) (*models.ControlProgramInfoRequest, error)

	StartPollingMachine(req connection.StartPollingRequest) *errors.AppError
	StopPollingMachine(machineID uint) *errors.AppError

	GetAllPollData() ([]entities.PollData, *errors.AppError)
	GetPollDataByID(recordID uint) (entities.PollData, *errors.AppError)
	SavePollData(data entities.PollData) *errors.AppError

	DeletePollDataByID(recordID uint) *errors.AppError

	// Получение batch необработанных данных
	GetPollDataBatch(batchSize int) ([]entities.PollData, *errors.AppError)

	// Пометка batch данных как обработанных
	MarkPollDataProcessed(ids []uint) *errors.AppError

	// Получение batch ID обработанных данных для удаления
	GetProcessedPollDataIDs(batchSize int) ([]uint, *errors.AppError)

	// Удаление batch данных по ID
	DeletePollDataBatch(ids []uint) *errors.AppError
}

package connection_usecase

import (
	"fmt"
	"log"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/pkg/errors"
	"strings"
	"time"
)

// validateAnonymousRequest проверяет обязательные поля для подключения анонимно
func (u *ConnectionUsecase) validateAnonymousRequest(request models.ConnectionRequest) error {
	if strings.TrimSpace(request.EndpointURL) == "" {
		return fmt.Errorf("endpoint URL is required")
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------------------

// ConnectAnonymous - анонимное подключение
func (u *ConnectionUsecase) ConnectAnonymous(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError) {
	var empty models.UUIDResponse

	if err := u.validateAnonymousRequest(request); err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "validation failed", err, true)
	}

	connReq := NewAnonymousConnectionFromRequest(&request)

	// Проверка доступности endpoint
	if err := isEndpointReachable(connReq.EndpointURL, 5*time.Second); err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "endpoint is not reachable", err, false)
	}

	if err := u.handleExistingMachine(connReq.EndpointURL); err != nil {
		return empty, err
	}

	machineUUID, err := u.createNewAnonymousConnection(connReq)
	if err != nil {
		return empty, err
	}

	log.Printf("✅ Successfully connected with UUID: %s", machineUUID)
	return models.UUIDResponse{UUID: machineUUID}, nil
}

// createNewAnonymousConnection создает анонимное соединение в сервисе и записи в БД
func (u *ConnectionUsecase) createNewAnonymousConnection(connReq *connection_models.AnonymousConnection) (string, *errors.AppError) {

	connID, err := u.OpcService.CreateAnonymousConnection(*connReq)
	if err != nil {
		return "", errors.NewAppError(errors.InternalServerErrorCode, "failed to create anonymous connection for machine", err, false)
	}

	newAnon := entities.AnonymousConnection{
		Policy: connReq.EndpointURL,
		Mode:   connReq.Manufacturer,
	}
	anonID, eerr := u.CreateAnonRecord(newAnon)
	if eerr != nil {
		return "", eerr
	}

	newMachine := entities.CncMachine{
		UUID:                  connID.String(),
		EndpointURL:           connReq.EndpointURL,
		Model:                 connReq.Model,
		Manufacturer:          connReq.Manufacturer,
		Status:                connection_models.ConnectionStatusConnected,
		Interval:              int(connReq.Timeout.Seconds()),
		ConnectionType:        connection_models.ConnectionAnonymous,
		AnonymousConnectionID: &anonID,
	}
	machineUUID, eerr := u.CreateMachineRecord(newMachine)
	if eerr != nil {
		return "", eerr
	}

	return machineUUID, nil
}

// ----------------------------------------------------------------------------------------------------------------

// NewAnonymousConnectionFromRequest Конструктор из ConnectionRequest
func NewAnonymousConnectionFromRequest(req *models.ConnectionRequest) *connection_models.AnonymousConnection {
	return &connection_models.AnonymousConnection{
		EndpointURL:  req.EndpointURL,
		Timeout:      time.Duration(req.Timeout) * time.Second,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}
}

package connection_usecase

import (
	"fmt"
	"log"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models/connection"
	connection_models "opc_ua_service/internal/domain/models/connection_types"
	"opc_ua_service/pkg/errors"
	"strings"
	"time"
)

// ConnectAnonymous - анонимное подключение
func (u *ConnectionUsecase) ConnectAnonymous(request connection.ConnectionRequest) (entities.CncMachine, *errors.AppError) {
	var empty entities.CncMachine

	if err := u.validateAnonymousRequest(request); err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "validation failed", err, true)
	}

	config := newAnonymousConnectionFromRequest(&request)

	// Проверка доступности endpoint
	if err := isEndpointReachable(config.EndpointURL, 5*time.Second); err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "endpoint is not reachable", err, false)
	}

	// Проверка существующей машины
	cfg := connection_models.ConnectionConfig{
		Config: config,
	}
	if err := u.handleExistingMachine(cfg); err != nil {
		return empty, err
	}

	machine, err := u.createNewAnonymousConnection(config)
	if err != nil {
		return empty, err
	}

	log.Printf("Successfully connected with ID: %s", machine.ID)
	return machine, nil
}

// createNewAnonymousConnection создает анонимное соединение в сервисе и записи в БД
func (u *ConnectionUsecase) createNewAnonymousConnection(connReq *connection_models.AnonymousConnection) (entities.CncMachine, *errors.AppError) {
	empty := entities.CncMachine{}

	connID, err := u.OpcService.CreateAnonymousConnection(*connReq)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to create anonymous connection for machine", err, false)
	}

	newAnon := entities.AnonymousConnection{
		Policy: connReq.EndpointURL,
		Mode:   connReq.Manufacturer,
	}
	anonID, eerr := u.CreateAnonRecord(newAnon)
	if eerr != nil {
		return empty, eerr
	}

	newMachine := entities.CncMachine{
		UUID:                  connID.String(),
		EndpointURL:           connReq.EndpointURL,
		Model:                 connReq.Model,
		Manufacturer:          connReq.Manufacturer,
		Status:                connection_models.ConnectionStatusConnected,
		PollTimeout:           int(connReq.Timeout.Seconds()),
		ConnectionType:        connection_models.ConnectionAnonymous,
		AnonymousConnectionID: &anonID,
	}
	machine, eerr := u.CreateMachineRecord(newMachine)
	if eerr != nil {
		return empty, eerr
	}

	return machine, nil
}

// ----------------------------------------------------------------------------------------------------------------

// newAnonymousConnectionFromRequest Конструктор из ConnectionRequest
func newAnonymousConnectionFromRequest(req *connection.ConnectionRequest) *connection_models.AnonymousConnection {
	return &connection_models.AnonymousConnection{
		EndpointURL:  req.EndpointURL,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}
}

// ----------------------------------------------------------------------------------------------------------------

// validateAnonymousRequest проверяет обязательные поля для подключения анонимно
func (u *ConnectionUsecase) validateAnonymousRequest(request connection.ConnectionRequest) error {
	if strings.TrimSpace(request.EndpointURL) == "" {
		return fmt.Errorf("endpoint URL is required")
	}
	return nil
}

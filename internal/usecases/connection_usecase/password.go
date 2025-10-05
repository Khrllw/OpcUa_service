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

// ConnectWithPassword - подключение по паролю
// ConnectWithPassword выполняет подключение по паролю и создает записи в БД
func (u *ConnectionUsecase) ConnectWithPassword(request connection.ConnectionRequest) (entities.CncMachine, *errors.AppError) {
	var empty entities.CncMachine

	if err := u.validatePasswordRequest(request); err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "validation failed", err, true)
	}

	config := newPasswordConnectionFromRequest(&request)

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

	machine, err := u.createNewPasswordConnection(config)
	if err != nil {
		return empty, err
	}

	log.Printf("Successfully connected with UUID: %s", machine.ID)
	return machine, nil
}

// createNewPasswordConnection создает соединение по паролю и записи в БД
func (u *ConnectionUsecase) createNewPasswordConnection(connReq *connection_models.PasswordConnection) (entities.CncMachine, *errors.AppError) {
	empty := entities.CncMachine{}
	connID, err := u.OpcService.CreatePasswordConnection(*connReq)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to create password connection for machine", err, false)
	}

	newPass := entities.PasswordConnection{
		Username: connReq.Username,
		Password: connReq.Password,
		Policy:   connReq.Policy,
		Mode:     connReq.Mode,
	}
	passID, eerr := u.CreatePasswordRecord(newPass)
	if eerr != nil {
		return empty, eerr
	}

	newMachine := entities.CncMachine{
		UUID:                 connID.String(),
		EndpointURL:          connReq.EndpointURL,
		Model:                connReq.Model,
		Manufacturer:         connReq.Manufacturer,
		Status:               connection_models.ConnectionStatusConnected,
		PollTimeout:          int(connReq.Timeout.Seconds()),
		ConnectionType:       connection_models.ConnectionPassword,
		PasswordConnectionID: &passID,
	}
	machine, eerr := u.CreateMachineRecord(newMachine)
	if eerr != nil {
		return empty, eerr
	}

	return machine, nil
}

// ----------------------------------------------------------------------------------------------------------------

// newPasswordConnectionFromRequest Конструктор из ConnectionRequest
func newPasswordConnectionFromRequest(req *connection.ConnectionRequest) *connection_models.PasswordConnection {
	return &connection_models.PasswordConnection{
		EndpointURL:  req.EndpointURL,
		Username:     req.Username,
		Password:     req.Password,
		Policy:       string(req.Policy),
		Mode:         string(req.Mode),
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}
}

// ----------------------------------------------------------------------------------------------------------------

// validatePasswordRequest проверяет обязательные поля для подключения по паролю
func (u *ConnectionUsecase) validatePasswordRequest(request connection.ConnectionRequest) error {
	if strings.TrimSpace(request.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(request.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

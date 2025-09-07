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

// validatePasswordRequest проверяет обязательные поля для подключения по паролю
func (u *ConnectionUsecase) validatePasswordRequest(request models.ConnectionRequest) error {
	if strings.TrimSpace(request.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(request.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------------------

// ConnectWithPassword - подключение по паролю
// ConnectWithPassword выполняет подключение по паролю и создает записи в БД
func (u *ConnectionUsecase) ConnectWithPassword(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError) {
	var empty models.UUIDResponse

	if err := u.validatePasswordRequest(request); err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "validation failed", err, true)
	}

	connReq := NewPasswordConnectionFromRequest(&request)

	// Проверка доступности endpoint
	if err := isEndpointReachable(connReq.EndpointURL, 5*time.Second); err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "endpoint is not reachable", err, false)
	}

	if err := u.handleExistingMachine(connReq.EndpointURL); err != nil {
		return empty, err
	}

	machineUUID, err := u.createNewPasswordConnection(connReq)
	if err != nil {
		return empty, err
	}

	log.Printf("✅ Successfully connected with UUID: %s", machineUUID)
	return models.UUIDResponse{UUID: machineUUID}, nil
}

// createNewPasswordConnection создает соединение по паролю и записи в БД
func (u *ConnectionUsecase) createNewPasswordConnection(connReq *connection_models.PasswordConnection) (string, *errors.AppError) {
	connID, err := u.OpcService.CreatePasswordConnection(*connReq)
	if err != nil {
		return "", errors.NewAppError(errors.InternalServerErrorCode, "failed to create password connection for machine", err, false)
	}

	newPass := entities.PasswordConnection{
		Username: connReq.Username,
		Password: connReq.Password,
		Policy:   connReq.Policy,
		Mode:     connReq.Mode,
	}
	passID, eerr := u.CreatePasswordRecord(newPass)
	if eerr != nil {
		return "", eerr
	}

	newMachine := entities.CncMachine{
		UUID:                 connID.String(),
		EndpointURL:          connReq.EndpointURL,
		Model:                connReq.Model,
		Manufacturer:         connReq.Manufacturer,
		Status:               connection_models.ConnectionStatusConnected,
		Interval:             int(connReq.Timeout.Seconds()),
		ConnectionType:       connection_models.ConnectionPassword,
		PasswordConnectionID: &passID,
	}
	machineUUID, eerr := u.CreateMachineRecord(newMachine)
	if eerr != nil {
		return "", eerr
	}

	return machineUUID, nil
}

// ----------------------------------------------------------------------------------------------------------------

// NewPasswordConnectionFromRequest Конструктор из ConnectionRequest
func NewPasswordConnectionFromRequest(req *models.ConnectionRequest) *connection_models.PasswordConnection {
	return &connection_models.PasswordConnection{
		EndpointURL:  req.EndpointURL,
		Username:     req.Username,
		Password:     req.Password,
		Policy:       string(req.Policy),
		Mode:         string(req.Mode),
		Timeout:      time.Duration(req.Timeout) * time.Second,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}
}

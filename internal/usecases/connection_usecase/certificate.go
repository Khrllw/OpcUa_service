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

// validateCertificateRequest проверяет обязательные поля для подключения по сертификату
func (u *ConnectionUsecase) validateCertificateRequest(request models.ConnectionRequest) error {
	if strings.TrimSpace(request.EndpointURL) == "" {
		return fmt.Errorf("endpoint URL is required")
	}
	if strings.TrimSpace(request.Certificate) == "" {
		return fmt.Errorf("certificate is required")
	}
	if strings.TrimSpace(request.Key) == "" {
		return fmt.Errorf("private key is required")
	}
	return nil
}

// ----------------------------------------------------------------------------------------------------------------

// ConnectWithCertificate выполняет подключение по сертификату и создает записи в БД
func (u *ConnectionUsecase) ConnectWithCertificate(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError) {
	var empty models.UUIDResponse

	// Валидация и подготовка структуры соединения
	if err := u.validateCertificateRequest(request); err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "validation failed", err, true)
	}

	connReq, err := u.NewCertificateConnectionFromRequest(&request)
	if err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "invalid request", err, true)
	}

	// Проверка доступности endpoint
	if err := isEndpointReachable(connReq.EndpointURL, 5*time.Second); err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "endpoint is not reachable", err, false)
	}

	config := connection_models.CertificateConnection{
		EndpointURL:  connReq.EndpointURL,
		Certificate:  connReq.Certificate,
		Key:          connReq.Key,
		Policy:       connReq.Policy,
		Mode:         connReq.Mode,
		Timeout:      time.Duration(request.Timeout) * time.Second,
		Manufacturer: connReq.Manufacturer,
		Model:        connReq.Model,
	}

	// Проверка существующей машины и закрытие старого соединения
	if err := u.handleExistingMachine(config.EndpointURL); err != nil {
		return empty, err
	}

	// Создание нового соединения и запись в БД
	machineUUID, eerr := u.createNewCertificateConnection(config)
	if eerr != nil {
		return empty, eerr
	}

	log.Printf("✅ Successfully connected with UUID: %s", machineUUID)
	return models.UUIDResponse{UUID: machineUUID}, nil
}

// createNewCertificateConnection создает соединение в сервисе и записи в БД
func (u *ConnectionUsecase) createNewCertificateConnection(config connection_models.CertificateConnection) (string, *errors.AppError) {

	connID, err := u.OpcService.CreateCertificateConnection(config)
	if err != nil {
		return "", errors.NewAppError(errors.InternalServerErrorCode, "failed to create connection for machine", err, false)
	}

	newCert := entities.CertificateConnection{
		Certificate: config.Certificate,
		Key:         config.Key,
		Policy:      config.Policy,
		Mode:        config.Mode,
	}
	certID, eerr := u.CreateCertRecord(newCert)
	if eerr != nil {
		return "", eerr
	}

	newMachine := entities.CncMachine{
		UUID:                    connID.String(),
		EndpointURL:             config.EndpointURL,
		Model:                   config.Model,
		Manufacturer:            config.Manufacturer,
		Status:                  connection_models.ConnectionStatusConnected,
		Interval:                int(config.Timeout.Seconds()),
		ConnectionType:          "certificate",
		CertificateConnectionID: &certID,
	}

	machineUUID, eerr := u.CreateMachineRecord(newMachine)
	if eerr != nil {
		return "", eerr
	}

	return machineUUID, nil
}

// ----------------------------------------------------------------------------------------------------------------

// NewCertificateConnectionFromRequest Конструктор из ConnectionRequest
func (u *ConnectionUsecase) NewCertificateConnectionFromRequest(req *models.ConnectionRequest) (connection_models.CertificateConnection, error) {
	var empty connection_models.CertificateConnection

	// Попытка распарсить сертификат (Base64 -> []byte)
	parsedCert, err := u.OpcService.Base64ToBytes(req.Certificate)
	if err != nil {
		return empty, err
	}

	// Попытка распарсить ключ (Base64 -> []byte)
	parsedKey, err := u.OpcService.Base64ToBytes(req.Key)
	if err != nil {
		return empty, err
	}

	// Валидация: всё же убедимся, что ключ и сертификат можно распарсить
	_, cert, key := u.OpcService.DecodeClientCredentials(parsedCert, parsedKey)
	if cert == nil || key == nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "invalid certificate or key content", nil, false)
	}

	// Успешное создание подключения
	return connection_models.CertificateConnection{
		EndpointURL:  req.EndpointURL,
		Certificate:  parsedCert,
		Key:          parsedKey,
		Policy:       string(req.Policy),
		Mode:         string(req.Mode),
		Timeout:      time.Duration(req.Timeout) * time.Second,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}, nil
}

package usecases

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/pkg/errors"
	"strings"
	"time"
)

type ConnectionUsecase struct {
	OpcService   interfaces.OpcService
	MachineRepo  interfaces.CncMachineRepository
	CertRepo     interfaces.CertificateConnectionRepository
	PasswordRepo interfaces.PasswordConnectionRepository
	AnonRepo     interfaces.AnonymousConnectionRepository
}

func NewConnectionUsecase(s interfaces.OpcService, r interfaces.CncMachineRepository, cr interfaces.CertificateConnectionRepository, pr interfaces.PasswordConnectionRepository, ar interfaces.AnonymousConnectionRepository) *ConnectionUsecase {
	return &ConnectionUsecase{s, r, cr, pr, ar}
}

// ConnectAnonymous - анонимное подключение
func (u *ConnectionUsecase) ConnectAnonymous(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError) {
	var empty models.UUIDResponse
	NewAnonymousConnectionFromRequest(&request)
	return empty, nil
}

// ConnectWithPassword - подключение по паролю
func (u *ConnectionUsecase) ConnectWithPassword(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError) {
	var empty models.UUIDResponse
	if err := u.validatePasswordRequest(request); err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "validation failed", err, true)
	}
	NewPasswordConnectionFromRequest(&request)
	return empty, nil
}

// ConnectWithCertificate - подключение по сертификату
func (u *ConnectionUsecase) ConnectWithCertificate(request models.ConnectionRequest) (models.UUIDResponse, *errors.AppError) {
	var empty models.UUIDResponse

	if err := u.validateCertificateRequest(request); err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "validation failed", err, true)
	}

	req, err := u.NewCertificateConnectionFromRequest(&request)
	if err != nil {
		return empty, errors.NewAppError(errors.InvalidDataCode, "invalid request", err, true)
	}

	config := connection_models.CertificateConnection{
		EndpointURL:  req.EndpointURL,
		Certificate:  req.Certificate,
		Key:          req.Key,
		Policy:       req.Policy,
		Mode:         req.Mode,
		Timeout:      time.Duration(request.Timeout) * time.Second,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}

	machineUUID, eerr := u.сreateConnection(config)
	if eerr != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to create connection", err, false)
	}

	log.Printf("Successfully connected with UUID: %s", machineUUID)
	return models.UUIDResponse{UUID: machineUUID}, nil
}

// сreateConnection проверяет наличие машины в БД и активного соединения в пуле.
// Если соединение отсутствует — оно будет создано.
// Возвращает UUID машины.
func (u *ConnectionUsecase) сreateConnection(req connection_models.CertificateConnection) (string, *errors.AppError) {
	var empty = ""

	// Проверяем наличие машины по EndpointURL
	foundMachine, err := u.MachineRepo.GetCncMachineByEndpointURL(req.EndpointURL)
	fmt.Println(errors.Is(err, errors.ErrNotFound))
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, errors.InternalServerError, err, false)
	}

	// Если машина есть, проверим, есть ли соединение и корректно его закроем
	if foundMachine.UUID != "" {
		id, err := uuid.Parse(foundMachine.UUID)
		if err != nil {
			return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to parse exist machine UUID", err, false)
		}
		conn, err := u.OpcService.GetConnectionByUUID(id)
		if conn != nil {
			if _, eerr := u.DisconnectByUUID(id); eerr != nil {
				if !errors.Is(eerr, errors.ErrNotFound) {
					return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to disconnect old machine", err, false)
				}
			}
		} else {
			if err := u.MachineRepo.DeleteCncMachine(id.String()); err != nil {
				if !errors.Is(err, errors.ErrNotFound) {
					return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to delete old machine record", err, false)
				}
			}
			/*
				if foundMachine.ConnectionType == "certificate" {
					certInfo := *foundMachine.CertificateConnectionID
					eerr := u.CertRepo.DeleteCertificateConnection(certInfo)
					if !errors.Is(eerr, errors.ErrNotFound) {
						return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to disconnect old machine", err, false)
					}
				}

			*/

		}
	}

	// Соединения нет — создаем новое
	connID, err := u.OpcService.CreateCertificateConnection(req)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to create connection for machine", err, false)
	}

	// Добавляем в БД
	newCertInfo := entities.CertificateConnection{
		Certificate: req.Certificate,
		Key:         req.Key,
		Policy:      req.Policy,
		Mode:        req.Mode,
	}

	certInfoID, err := u.CertRepo.CreateCertificateConnection(newCertInfo)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to create certificate info record", err, false)
	}

	newMachine := entities.CncMachine{
		UUID:                    connID.String(),
		EndpointURL:             req.EndpointURL,
		Model:                   req.Model,
		Manufacturer:            req.Manufacturer,
		Status:                  connection_models.ConnectionStatusConnected,
		Interval:                int(req.Timeout.Seconds()),
		ConnectionType:          "certificate",
		CertificateConnectionID: &certInfoID,
	}

	machineUUID, err := u.MachineRepo.CreateCncMachine(newMachine)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to create machine record", err, false)
	}

	// Проверяем, что машина добавлена
	addedMachine, err := u.MachineRepo.GetCncMachineByUUID(machineUUID)
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to check machine connection", err, false)
	}

	return addedMachine.UUID, nil
}

// ----------------------------------------------------------------------------------------------------------------

// DisconnectByUUID закрывает соединение по UUID
func (u *ConnectionUsecase) DisconnectByUUID(id uuid.UUID) (*bool, *errors.AppError) {
	var state = false
	info, err := u.OpcService.GetConnectionInfoByUUID(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.NewAppError(errors.NotFoundErrorCode, "failed to find connection", err, false)
		} else {
			return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to get connection", err, false)
		}
	}
	if info.IsPolled {
		if err := u.OpcService.StopPollingForMachine(id); err != nil {
			log.Printf("Warning: failed to stop polling for machine %s: %v", id, err)
		}
	}

	if err := u.OpcService.CloseConnection(id); err != nil {
		return &state, errors.NewAppError(errors.InternalServerErrorCode, "failed to close connection", err, false)
	}
	state = true

	machine, err := u.MachineRepo.GetCncMachineByUUID(id.String())
	if err != nil {
		return &state, errors.NewAppError(errors.InternalServerErrorCode, "failed to delete machine record", err, false)
	}
	if machine.ConnectionType == "certificate" {
		err = u.CertRepo.DeleteCertificateConnection(*machine.CertificateConnectionID)
		if err != nil {
			if !errors.Is(err, errors.ErrNotFound) {
				return &state, errors.NewAppError(errors.InternalServerErrorCode, "failed to delete certificate record", err, false)
			}

		}
	}

	err = u.MachineRepo.DeleteCncMachine(id.String())
	if err != nil {
		if !errors.Is(err, errors.ErrNotFound) {
			return &state, errors.NewAppError(errors.InternalServerErrorCode, "failed to delete machine record", err, false)
		}
	}

	log.Printf("Successfully closed connection with UUID: %s", id)
	return &state, nil
}

// DisconnectAll закрывает все соединения
func (u *ConnectionUsecase) DisconnectAll() (int, *errors.AppError) {
	var empty = 0

	stats := u.OpcService.GetGlobalStats()
	activeConnections := stats.ActiveConnections

	u.OpcService.CloseAll()

	err := u.MachineRepo.DeleteAllCncMachines()
	if err != nil {
		return empty, errors.NewAppError(errors.InternalServerErrorCode, "failed to delete machine records", err, false)
	}

	log.Printf("Closed all %d active connections", activeConnections)
	return int(activeConnections), nil
}

// CleanupIdleConnections очищает неиспользуемые соединения
func (u *ConnectionUsecase) CleanupIdleConnections(maxIdleMinutes int) int {
	cleaned := u.OpcService.Cleanup(time.Duration(maxIdleMinutes) * time.Minute)
	log.Printf("Cleaned up %d idle connections (idle time > %d minutes)", cleaned, maxIdleMinutes)
	return cleaned
}

// ----------------------------------------------------------------------------------------------------------------

// GetActiveConnections возвращает список активных соединений
func (u *ConnectionUsecase) GetActiveConnections() models.ConnectionPoolResponse {
	// Получаем все соединения из сервиса
	connectionsInfo := u.OpcService.GetAllConnectionsInfo()

	// Преобразуем каждый ConnectionInfo в ConnectionInfoResponse
	var result []*models.ConnectionInfoResponse
	for id, connInfo := range connectionsInfo {
		response := u.сonvertConnectionInfoToResponse(id, connInfo)
		result = append(result, &response)
	}

	return models.ConnectionPoolResponse{
		PoolSize:    len(result),
		Connections: result,
	}
}

// GetConnectionState получает состояние подключения
func (u *ConnectionUsecase) GetConnectionState(id uuid.UUID) (*models.ConnectionInfoResponse, *errors.AppError) {
	connInfo, err := u.OpcService.GetConnectionInfoByUUID(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, errors.NewAppError(errors.NotFoundErrorCode, "failed to find connection", errors.ErrNotFound, false)
		} else {
			return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to get connection", err, false)
		}
	}
	stat := u.сonvertConnectionInfoToResponse(id, connInfo)

	return &stat, nil
}

// ----------------------------------------------------------------------------------------------------------------

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

// validateAnonymousRequest проверяет обязательные поля для подключения анонимно
func (u *ConnectionUsecase) validateAnonymousRequest(request models.ConnectionRequest) error {
	if strings.TrimSpace(request.EndpointURL) == "" {
		return fmt.Errorf("endpoint URL is required")
	}
	return nil
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
	_, cert, key := u.OpcService.LoadClientCredentialsBytes(parsedCert, parsedKey)
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

// NewAnonymousConnectionFromRequest Конструктор из ConnectionRequest
func NewAnonymousConnectionFromRequest(req *models.ConnectionRequest) *connection_models.AnonymousConnection {
	return &connection_models.AnonymousConnection{
		EndpointURL:  req.EndpointURL,
		Timeout:      time.Duration(req.Timeout) * time.Second,
		Manufacturer: req.Manufacturer,
		Model:        req.Model,
	}
}

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

// сonvertConnectionInfoToResponse преобразует ConnectionInfo в ConnectionInfoResponse
func (u *ConnectionUsecase) сonvertConnectionInfoToResponse(id uuid.UUID, connInfo *models.ConnectionInfo) models.ConnectionInfoResponse {
	if connInfo == nil {
		return models.ConnectionInfoResponse{
			Status:      models.StatusNotFound,
			Description: models.StatusNotFound.GetDescription(),
		}
	}

	// Блокируем для безопасного чтения
	connInfo.Mu.RLock()
	defer connInfo.Mu.RUnlock()

	// Определяем статус на основе IsHealthy
	var status models.ConnectionStatusEnum
	if connInfo.IsHealthy {
		status = models.StatusHealthy
	} else {
		status = models.StatusUnhealthy
	}

	return models.ConnectionInfoResponse{
		UUID:        id,
		SessionID:   connInfo.SessionID,
		Status:      status,
		Description: status.GetDescription(),
		Config:      connInfo.Config,
		CreatedAt:   connInfo.CreatedAt,
		LastUsed:    connInfo.LastUsed,
		UseCount:    connInfo.UseCount,
	}
}

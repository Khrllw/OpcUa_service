package connection_usecase

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"net"
	"net/url"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/internal/interfaces"
	"opc_ua_service/pkg/errors"
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
	eerr := u.DeleteMachineRecord(id)
	if eerr != nil {
		return &state, eerr
	}

	log.Printf("Successfully closed connection with UUID: %s", id)
	return &state, nil
}

// ----------------------------------------------------------------------------------------------------------------

// handleExistingMachine проверяет наличие машины по EndpointURL и закрывает старое соединение
func (u *ConnectionUsecase) handleExistingMachine(endpointURL string) *errors.AppError {
	machine, err := u.MachineRepo.GetCncMachineByEndpointURL(endpointURL)
	if err != nil && !errors.Is(err, errors.ErrNotFound) {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to get machine by endpoint", err, false)
	}

	if machine.UUID == "" {
		return nil
	}

	id, err := uuid.Parse(machine.UUID)
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to parse existing machine UUID", err, false)
	}

	conn, _ := u.OpcService.GetConnectionByUUID(id)
	if conn != nil {
		if _, err := u.DisconnectByUUID(id); err != nil && !errors.Is(err, errors.ErrNotFound) {
			return errors.NewAppError(errors.InternalServerErrorCode, "failed to disconnect old machine", err, false)
		}
	} else {
		if err := u.DeleteMachineRecord(id); err != nil && !errors.Is(err, errors.ErrNotFound) {
			return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete old machine record", err, false)
		}
	}

	return nil
}

// ----------------------------------------------------------------------------------------------------------------

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

// isEndpointReachable проверяет, доступен ли TCP порт указанного endpoint за заданный таймаут.
func isEndpointReachable(endpoint string, timeout time.Duration) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("invalid endpoint URL: %w", err)
	}

	hostPort := fmt.Sprintf("%s:%s", u.Hostname(), u.Port())
	conn, err := net.DialTimeout("tcp", hostPort, timeout)
	if err != nil {
		return fmt.Errorf("endpoint %s is not reachable: %w", hostPort, err)
	}
	defer conn.Close()
	return nil
}

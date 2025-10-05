package connection_usecase

import (
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_types"
	"opc_ua_service/pkg/errors"
	"time"
)

// RestoreConnection восстанавливает подключение из БД в пул памяти.
func (u *ConnectionUsecase) RestoreConnection(machine entities.CncMachine) (*models.ConnectionInfo, *errors.AppError) {
	var connID *uuid.UUID
	var err error

	ID, err := uuid.Parse(machine.UUID)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to parse UUID", err, true)
	}

	switch machine.ConnectionType {
	case connection_models.ConnectionCertificate:
		req := connection_models.CertificateConnection{
			EndpointURL:  machine.EndpointURL,
			Model:        machine.Model,
			Manufacturer: machine.Manufacturer,
			Timeout:      time.Duration(machine.PollTimeout) * time.Second,
			Policy:       machine.CertificateConnection.Policy,
			Mode:         machine.CertificateConnection.Mode,
			Certificate:  machine.CertificateConnection.Certificate,
			Key:          machine.CertificateConnection.Key,
		}
		connID, err = u.OpcService.CreateCertificateConnection(req, &ID)

	case connection_models.ConnectionPassword:
		req := connection_models.PasswordConnection{
			EndpointURL:  machine.EndpointURL,
			Model:        machine.Model,
			Manufacturer: machine.Manufacturer,
			Timeout:      time.Duration(machine.PollTimeout) * time.Second,
			Username:     machine.PasswordConnection.Username,
			Password:     machine.PasswordConnection.Password,
			Policy:       machine.PasswordConnection.Policy,
			Mode:         machine.PasswordConnection.Mode,
		}
		connID, err = u.OpcService.CreatePasswordConnection(req)

	case connection_models.ConnectionAnonymous:
		req := connection_models.AnonymousConnection{
			EndpointURL:  machine.EndpointURL,
			Model:        machine.Model,
			Manufacturer: machine.Manufacturer,
			Timeout:      time.Duration(machine.PollTimeout) * time.Second,
		}
		connID, err = u.OpcService.CreateAnonymousConnection(req)

	default:
		return nil, errors.NewAppError(errors.InvalidDataCode, "invalid connection type", nil, true)
	}

	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to create connection", err, true)
	}
	if connID == nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to create connection", nil, true)
	}

	// Получаем информацию о соединении из пула
	connInfo, err2 := u.OpcService.GetConnectionInfoByUUID(*connID)
	if err2 != nil || connInfo == nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to add machine to pool", err, true) // соединение не удалось восстановить, но функция всегда успешна
	}
	connInfo.MachineID = machine.ID

	// Запускаем опрос, если машина была в состоянии "polled"
	if machine.Status == connection_models.ConnectionStatusPolled {
		if err := u.OpcService.StartPollingForMachine(*connID, connInfo.Config.Config.GetTimeout()); err != nil {
			return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to start polling for machine", err, true)
		}
	}

	return connInfo, nil
}

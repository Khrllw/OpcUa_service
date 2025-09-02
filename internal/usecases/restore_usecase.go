package usecases

import (
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/pkg/errors"
	"time"
)

// RestoreConnection восстанавливает подключение из БД в пул памяти.
func (u *ConnectionUsecase) RestoreConnection(machine entities.CncMachine) (*models.ConnectionInfo, *errors.AppError) {
	var connID *uuid.UUID
	var err error

	switch machine.ConnectionType {
	case connection_models.ConnectionCertificate:
		req := connection_models.CertificateConnection{
			EndpointURL:  machine.EndpointURL,
			Model:        machine.Model,
			Manufacturer: machine.Manufacturer,
			Timeout:      time.Duration(machine.Interval) * time.Second,
			Policy:       machine.CertificateConnection.Policy,
			Mode:         machine.CertificateConnection.Mode,
			Certificate:  machine.CertificateConnection.Certificate,
			Key:          machine.CertificateConnection.Key,
		}
		connID, err = u.OpcService.CreateCertificateConnection(req)

	case connection_models.ConnectionPassword:
		req := connection_models.PasswordConnection{
			EndpointURL:  machine.EndpointURL,
			Model:        machine.Model,
			Manufacturer: machine.Manufacturer,
			Timeout:      time.Duration(machine.Interval) * time.Second,
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
			Timeout:      time.Duration(machine.Interval) * time.Second,
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

	// Обновляем UUID машины в БД
	updateMap := map[string]interface{}{
		"UUID": connID.String(),
	}
	if _, eerr := u.MachineRepo.UpdateCncMachine(machine.UUID, updateMap); eerr != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to update machine UUID", err, true)
	}

	// Получаем информацию о соединении из пула
	connInfo, err2 := u.OpcService.GetConnectionInfoByUUID(*connID)
	if err2 != nil || connInfo == nil {
		return nil, nil // соединение не удалось восстановить, но функция всегда успешна
	}

	// Запускаем опрос, если машина была в состоянии "polled"
	if machine.Status == connection_models.ConnectionStatusPolled {
		if err := u.OpcService.StartPollingForMachine(*connID); err != nil {
			return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to start polling for machine", err, true)
		}
	}

	return connInfo, nil
}

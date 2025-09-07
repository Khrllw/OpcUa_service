package connection_usecase

import (
	"github.com/google/uuid"
	"opc_ua_service/internal/domain/entities"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/pkg/errors"
)

// CreateMachineRecord создает запись в БД
func (u *ConnectionUsecase) CreateMachineRecord(newMachine entities.CncMachine) (string, *errors.AppError) {
	var empty = ""
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

// CreateCertRecord создает запись в БД
func (u *ConnectionUsecase) CreateCertRecord(newCert entities.CertificateConnection) (uint, *errors.AppError) {
	createdCertID, err := u.CertRepo.CreateCertificateConnection(newCert)
	if err != nil {
		return 0, errors.NewAppError(errors.InternalServerErrorCode, "failed to create machine record", err, false)
	}

	addedCert, err := u.CertRepo.GetCertificateConnectionByID(createdCertID)
	if err != nil {
		return 0, errors.NewAppError(errors.InternalServerErrorCode, "failed to check machine connection", err, false)
	}
	return addedCert.ID, nil
}

// CreatePasswordRecord создает запись в БД
func (u *ConnectionUsecase) CreatePasswordRecord(newPass entities.PasswordConnection) (uint, *errors.AppError) {
	createdID, err := u.PasswordRepo.CreatePasswordConnection(newPass)
	if err != nil {
		return 0, errors.NewAppError(errors.InternalServerErrorCode, "failed to create password record", err, false)
	}

	addedPass, err := u.PasswordRepo.GetPasswordConnectionByID(createdID)
	if err != nil {
		return 0, errors.NewAppError(errors.InternalServerErrorCode, "failed to check password record", err, false)
	}
	return addedPass.ID, nil
}

// CreateAnonRecord создает запись в БД
func (u *ConnectionUsecase) CreateAnonRecord(newAnon entities.AnonymousConnection) (uint, *errors.AppError) {
	createdID, err := u.AnonRepo.CreateAnonymousConnection(newAnon)
	if err != nil {
		return 0, errors.NewAppError(errors.InternalServerErrorCode, "failed to create anonymous record", err, false)
	}

	addedAnon, err := u.AnonRepo.GetAnonymousConnectionByID(createdID)
	if err != nil {
		return 0, errors.NewAppError(errors.InternalServerErrorCode, "failed to check anonymous record", err, false)
	}
	return addedAnon.ID, nil
}

// DeleteMachineRecord удаляет запись в БД
func (u *ConnectionUsecase) DeleteMachineRecord(id uuid.UUID) *errors.AppError {
	machine, err := u.MachineRepo.GetCncMachineByUUID(id.String())
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete machine record", err, false)
	}
	switch machine.ConnectionType {
	case connection_models.ConnectionCertificate:
		if machine.CertificateConnectionID != nil {
			err = u.CertRepo.DeleteCertificateConnection(*machine.CertificateConnectionID)
			if err != nil && !errors.Is(err, errors.ErrNotFound) {
				return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete certificate record", err, false)
			}
		}
	case connection_models.ConnectionPassword:
		if machine.PasswordConnectionID != nil {
			err = u.PasswordRepo.DeletePasswordConnection(*machine.PasswordConnectionID)
			if err != nil && !errors.Is(err, errors.ErrNotFound) {
				return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete password record", err, false)
			}
		}
	case connection_models.ConnectionAnonymous:
		if machine.AnonymousConnectionID != nil {
			err = u.AnonRepo.DeleteAnonymousConnection(*machine.AnonymousConnectionID)
			if err != nil && !errors.Is(err, errors.ErrNotFound) {
				return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete anonymous record", err, false)
			}
		}
	default:
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete record with unknown connection type", nil, false)
	}

	err = u.MachineRepo.DeleteCncMachine(id.String())
	if err != nil {
		if !errors.Is(err, errors.ErrNotFound) {
			return errors.NewAppError(errors.InternalServerErrorCode, "failed to delete machine record", err, false)
		}
	}
	return nil
}

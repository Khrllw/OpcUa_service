package cnc_machine

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/pkg/errors"
)

// CreateCncMachine создает новую запись о станке
func (r *CncMachineRepositoryImpl) CreateCncMachine(cnc *entities.CncMachine) (string, error) {
	op := "repo.CncMachine.CreateCncMachine"

	err := r.db.Clauses(clause.Returning{}).Create(&cnc).Error
	if err != nil {
		return "", errors.NewDBError(op, err)
	}

	return cnc.UUID, nil
}

// GetCncMachineByEndpointURL возвращает станок по адресу подключения
func (r *CncMachineRepositoryImpl) GetCncMachineByEndpointURL(endpoint string) (entities.CncMachine, error) {
	op := "repo.CncMachine.GetCncMachineByEndpointURL"

	var cnc entities.CncMachine
	err := r.db.First(&cnc, "endpoint_url = ?", endpoint).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.CncMachine{}, errors.NewDBError(op, errors.ErrNotFound)
		}
		return entities.CncMachine{}, errors.NewDBError(op, err)
	}

	return cnc, nil
}

// GetCncMachineByUUID возвращает станок по UUID
func (r *CncMachineRepositoryImpl) GetCncMachineByUUID(id string) (entities.CncMachine, error) {
	op := "repo.CncMachine.GetCncMachineByUUID"

	var cnc entities.CncMachine
	err := r.db.First(&cnc, "UUID = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.CncMachine{}, errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrNotFound))
		}
		return entities.CncMachine{}, errors.NewDBError(op, err)
	}

	return cnc, nil
}

// UpdateCncMachine обновляет поля станка по UUID
func (r *CncMachineRepositoryImpl) UpdateCncMachine(id string, updateMap map[string]interface{}) (string, error) {
	op := "repo.CncMachine.UpdateCncMachine"

	var updatedCnc entities.CncMachine
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedCnc).
		Where("UUID = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return "", errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return "", errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return updatedCnc.UUID, nil
}

// DeleteCncMachine удаляет станок по UUID
func (r *CncMachineRepositoryImpl) DeleteCncMachine(id string) error {
	op := "repo.CncMachine.DeleteCncMachine"

	result := r.db.Delete(&entities.CncMachine{}, "UUID = ?", id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}

// GetAllCncMachines возвращает список всех станков
func (r *CncMachineRepositoryImpl) GetAllCncMachines() ([]entities.CncMachine, error) {
	op := "repo.CncMachine.GetAllCncMachines"

	var list []entities.CncMachine
	if err := r.db.Find(&list).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return list, nil
}

// DeleteAllCncMachines удаляет все записи о станках
func (r *CncMachineRepositoryImpl) DeleteAllCncMachines() error {
	op := "repo.CncMachine.DeleteAllCncMachines"

	result := r.db.Where("1 = 1").Delete(&entities.CncMachine{})
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}

package cnc_machine

import (
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

	return cnc.SIK, nil
}

// GetCncMachineBySIK возвращает станок по SIK
func (r *CncMachineRepositoryImpl) GetCncMachineBySIK(sik string) (entities.CncMachine, error) {
	op := "repo.CncMachine.GetCncMachineBySIK"

	var cnc entities.CncMachine
	err := r.db.First(&cnc, "sik = ?", sik).Error
	if err != nil {
		return entities.CncMachine{}, errors.NewDBError(op, err)
	}

	return cnc, nil
}

// UpdateCncMachine обновляет поля станка по SIK
func (r *CncMachineRepositoryImpl) UpdateCncMachine(sik string, updateMap map[string]interface{}) (string, error) {
	op := "repo.CncMachine.UpdateCncMachine"

	var updatedCnc entities.CncMachine
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updatedCnc).
		Where("sik = ?", sik).
		Updates(updateMap)

	if result.Error != nil {
		return "", errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return "", errors.NewDBError(op, result.Error)
	}

	return updatedCnc.SIK, nil
}

// DeleteCncMachine удаляет станок по SIK
func (r *CncMachineRepositoryImpl) DeleteCncMachine(sik string) error {
	op := "repo.CncMachine.DeleteCncMachine"

	result := r.db.Delete(&entities.CncMachine{}, "sik = ?", sik)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, result.Error)
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

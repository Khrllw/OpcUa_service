package password_connection

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/pkg/errors"
)

func (r *PasswordConnectionRepositoryImpl) CreatePasswordConnection(pc entities.PasswordConnection) (uint, error) {
	op := "repo.PasswordConnection.CreatePasswordConnection"

	err := r.db.Clauses(clause.Returning{}).Create(&pc).Error
	if err != nil {
		return 0, errors.NewDBError(op, err)
	}

	return pc.ID, nil
}

func (r *PasswordConnectionRepositoryImpl) GetPasswordConnectionByID(id uint) (entities.PasswordConnection, error) {
	op := "repo.PasswordConnection.GetPasswordConnectionByID"

	var pc entities.PasswordConnection
	err := r.db.First(&pc, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.PasswordConnection{}, errors.NewDBError(op, errors.ErrNotFound)
		}
		return entities.PasswordConnection{}, errors.NewDBError(op, err)
	}

	return pc, nil
}

func (r *PasswordConnectionRepositoryImpl) UpdatePasswordConnection(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.PasswordConnection.UpdatePasswordConnection"

	var updated entities.PasswordConnection
	result := r.db.
		Clauses(clause.Returning{}).
		Model(&updated).
		Where("id = ?", id).
		Updates(updateMap)

	if result.Error != nil {
		return 0, errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return 0, errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return updated.ID, nil
}

func (r *PasswordConnectionRepositoryImpl) DeletePasswordConnection(id uint) error {
	op := "repo.PasswordConnection.DeletePasswordConnection"

	result := r.db.Delete(&entities.PasswordConnection{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}

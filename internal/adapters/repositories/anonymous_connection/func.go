package anonymous_connection

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/pkg/errors"
)

func (r *AnonymousConnectionRepositoryImpl) CreateAnonymousConnection(ac entities.AnonymousConnection) (uint, error) {
	op := "repo.AnonymousConnection.CreateAnonymousConnection"

	err := r.db.Clauses(clause.Returning{}).Create(&ac).Error
	if err != nil {
		return 0, errors.NewDBError(op, err)
	}

	return ac.ID, nil
}

func (r *AnonymousConnectionRepositoryImpl) GetAnonymousConnectionByID(id uint) (entities.AnonymousConnection, error) {
	op := "repo.AnonymousConnection.GetAnonymousConnectionByID"

	var ac entities.AnonymousConnection
	err := r.db.First(&ac, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.AnonymousConnection{}, errors.NewDBError(op, errors.ErrNotFound)
		}
		return entities.AnonymousConnection{}, errors.NewDBError(op, err)
	}

	return ac, nil
}

func (r *AnonymousConnectionRepositoryImpl) UpdateAnonymousConnection(id uint, updateMap map[string]interface{}) (uint, error) {
	op := "repo.AnonymousConnection.UpdateAnonymousConnection"

	var updated entities.AnonymousConnection
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

func (r *AnonymousConnectionRepositoryImpl) DeleteAnonymousConnection(id uint) error {
	op := "repo.AnonymousConnection.DeleteAnonymousConnection"

	result := r.db.Delete(&entities.AnonymousConnection{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}

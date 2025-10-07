package poll_data

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/pkg/errors"
)

func (r PollDataRepositoryImpl) CreatePollData(pd entities.PollData) (uint, error) {
	op := "repo.PollDataRepository.CreatePollData"

	err := r.db.Clauses(clause.Returning{}).Create(&pd).Error
	if err != nil {
		return 0, errors.NewDBError(op, err)
	}

	return pd.ID, nil
}

func (r PollDataRepositoryImpl) GetPollDataByID(id uint) (entities.PollData, error) {
	op := "repo.PollDataRepository.GetPollDataByID"
	var empty entities.PollData

	var pd entities.PollData
	err := r.db.First(&pd, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return empty, errors.NewDBError(op, errors.ErrNotFound)
		}
		return empty, errors.NewDBError(op, err)
	}

	return pd, nil
}

func (r PollDataRepositoryImpl) GetAllPollData() ([]entities.PollData, error) {
	op := "repo.PollDataRepository.GetAllPollData"

	var list []entities.PollData
	if err := r.db.Find(&list).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return list, nil
}

func (r PollDataRepositoryImpl) DeletePollDataByID(id uint) error {
	op := "repo.PollDataRepository.DeletePollDataByID"

	result := r.db.Delete(&entities.PollData{}, "id = ?", id)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}

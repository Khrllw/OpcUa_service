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

// GetPollDataBatch получает batch необработанных записей
func (r *PollDataRepositoryImpl) GetPollDataBatch(batchSize int) ([]entities.PollData, error) {
	op := "repo.PollDataRepository.GetPollDataBatch"

	var records []entities.PollData
	if err := r.db.
		Where("processed = ?", false).
		Order("id ASC").
		Limit(batchSize).
		Find(&records).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return records, nil
}

// MarkPollDataProcessed помечает записи как обработанные
func (r *PollDataRepositoryImpl) MarkPollDataProcessed(ids []uint) error {
	op := "repo.PollDataRepository.MarkPollDataProcessed"

	result := r.db.Model(&entities.PollData{}).
		Where("ID IN ? AND processed = ?", ids, false).
		Update("processed", true)
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}

// GetProcessedPollDataIDs получает batch обработанных записей для удаления
func (r *PollDataRepositoryImpl) GetProcessedPollDataIDs(batchSize int) ([]uint, error) {
	op := "repo.PollDataRepository.GetProcessedPollDataIDs"

	var ids []uint
	if err := r.db.
		Model(&entities.PollData{}).
		Where("processed = ?", true).
		Order("id ASC").
		Limit(batchSize).
		Pluck("id", &ids).Error; err != nil {
		return nil, errors.NewDBError(op, err)
	}

	return ids, nil
}

// DeletePollDataBatch удаляет записи по batch ID
func (r *PollDataRepositoryImpl) DeletePollDataBatch(ids []uint) error {
	op := "repo.PollDataRepository.DeletePollDataBatch"

	result := r.db.
		Where("id IN ?", ids).
		Delete(&entities.PollData{})
	if result.Error != nil {
		return errors.NewDBError(op, result.Error)
	}
	if result.RowsAffected == 0 {
		return errors.NewDBError(op, fmt.Errorf("%s: %w", op, errors.ErrEmptyAction))
	}

	return nil
}

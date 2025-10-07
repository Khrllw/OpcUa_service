package polling_usecase

import (
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/pkg/errors"
)

func (u *PollingUsecase) GetAllPollData() ([]entities.PollData, *errors.AppError) {
	data, err := u.PollDataRepo.GetAllPollData()
	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "", err, false)
	}
	return data, nil
}

func (u *PollingUsecase) GetPollDataByID(entryID uint) (entities.PollData, *errors.AppError) {
	var empty entities.PollData
	return empty, nil
}

func (u *PollingUsecase) SavePollData(data entities.PollData) *errors.AppError {
	return nil
}

func (u *PollingUsecase) DeletePollDataByID(entryID uint) *errors.AppError {
	if err := u.PollDataRepo.DeletePollDataByID(entryID); err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "", err, false)
	}
	return nil
}

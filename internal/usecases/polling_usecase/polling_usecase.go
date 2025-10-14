package polling_usecase

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"opc_ua_service/internal/domain/entities"
	"opc_ua_service/internal/domain/models"
	"opc_ua_service/internal/domain/models/connection"
	connection_models "opc_ua_service/internal/domain/models/connection_types"
	"opc_ua_service/pkg/errors"
	"time"

	"opc_ua_service/internal/interfaces"
)

type PollingUsecase struct {
	OpcService   interfaces.OpcService
	poolingCtx   context.Context
	cancelFunc   context.CancelFunc
	latestData   []interfaces.MachineData
	MachineRepo  interfaces.CncMachineRepository
	PollDataRepo interfaces.PollDataRepository
}

func NewPollingUsecase(s interfaces.OpcService, cnc_r interfaces.CncMachineRepository, poll_r interfaces.PollDataRepository) *PollingUsecase {
	return &PollingUsecase{
		OpcService:   s,
		MachineRepo:  cnc_r,
		PollDataRepo: poll_r,
	}
}

// GetControlProgram - получение управляющей программы
func (u *PollingUsecase) GetControlProgram(req connection.GetControlProgramRequest) (*models.ControlProgramInfoRequest, error) {
	machine, err := u.MachineRepo.GetCncMachineByID(req.ID)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to get machine", err, false)
	}

	id, err := uuid.Parse(machine.UUID)
	if err != nil {
		return nil, errors.NewAppError(errors.InternalServerErrorCode, "failed to parse UUID", err, false)
	}

	info, err := u.OpcService.GetControlProgramInfo(id)
	if err != nil {
		return nil, err
	}

	resp := &models.ControlProgramInfoRequest{
		ExecutionStack: info,
	}

	return resp, nil
}

// StartPollingMachine запускает сбор данных машины по UUID
func (u *PollingUsecase) StartPollingMachine(req connection.StartPollingRequest) *errors.AppError {
	machine, err := u.MachineRepo.GetCncMachineByID(req.ID)
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to get machine", err, false)
	}

	id, err := uuid.Parse(machine.UUID)
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to parse UUID", err, false)
	}

	poll_timeout := time.Duration(req.Timeout) * time.Second
	err = u.OpcService.StartPollingForMachine(id, poll_timeout)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return errors.NewAppError(http.StatusNotFound, "machine not found", err, false)
		} else {
			return errors.NewAppError(http.StatusInternalServerError, "failed to start polling for machine", err, false)
		}
	}
	updateMap := map[string]interface{}{
		"status":       connection_models.ConnectionStatusPolled,
		"poll_timeout": req.Timeout,
	}

	_, err = u.MachineRepo.UpdateCncMachine(machine.UUID, updateMap)
	if err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "failed to update machine record", err, false)
	}
	return nil
}

// StopPollingMachine останавливает сбор данных машины по UUID
func (u *PollingUsecase) StopPollingMachine(machineID uint) *errors.AppError {
	machine, err := u.MachineRepo.GetCncMachineByID(machineID)
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to get machine", err, false)
	}

	id, err := uuid.Parse(machine.UUID)
	if err != nil {
		return errors.NewAppError(errors.InternalServerErrorCode, "failed to parse UUID", err, false)
	}

	err = u.OpcService.StopPollingForMachine(id)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return errors.NewAppError(http.StatusNotFound, "machine not found", err, false)
		} else {
			return errors.NewAppError(http.StatusInternalServerError, "failed to stop polling for machine", err, false)
		}
	}
	updateMap := map[string]interface{}{
		"status": connection_models.ConnectionStatusConnected,
	}

	_, err = u.MachineRepo.UpdateCncMachine(machine.UUID, updateMap)
	if err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "failed to update machine record", err, false)
	}
	return nil
}

// GetPollDataBatch получает batch необработанных данных
func (u *PollingUsecase) GetPollDataBatch(batchSize int) ([]entities.PollData, *errors.AppError) {
	data, err := u.PollDataRepo.GetPollDataBatch(batchSize)
	if err != nil {
		return nil, errors.NewAppError(http.StatusInternalServerError, "failed to get poll data batch", err, false)
	}
	return data, nil
}

// MarkPollDataProcessed помечает batch данных как обработанные
func (u *PollingUsecase) MarkPollDataProcessed(ids []uint) *errors.AppError {
	if len(ids) == 0 {
		return nil
	}
	if err := u.PollDataRepo.MarkPollDataProcessed(ids); err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "failed to mark poll data as processed", err, false)
	}
	return nil
}

// GetProcessedPollDataIDs получает batch ID обработанных данных для удаления
func (u *PollingUsecase) GetProcessedPollDataIDs(batchSize int) ([]uint, *errors.AppError) {
	ids, err := u.PollDataRepo.GetProcessedPollDataIDs(batchSize)
	if err != nil {
		return nil, errors.NewAppError(http.StatusInternalServerError, "failed to get processed poll data IDs", err, false)
	}
	return ids, nil
}

// DeletePollDataBatch удаляет batch данных по ID
func (u *PollingUsecase) DeletePollDataBatch(ids []uint) *errors.AppError {
	if len(ids) == 0 {
		return nil
	}
	if err := u.PollDataRepo.DeletePollDataBatch(ids); err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "failed to delete poll data batch", err, false)
	}
	return nil
}

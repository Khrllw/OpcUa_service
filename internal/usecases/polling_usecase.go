package usecases

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"opc_ua_service/internal/domain/models"
	connection_models "opc_ua_service/internal/domain/models/connection_models"
	"opc_ua_service/pkg/errors"

	"opc_ua_service/internal/interfaces"
)

type PollingUsecase struct {
	OpcService interfaces.OpcService
	poolingCtx context.Context
	cancelFunc context.CancelFunc
	latestData []interfaces.MachineData
	Repo       interfaces.CncMachineRepository
}

func NewPollingUsecase(s interfaces.OpcService, r interfaces.CncMachineRepository) *PollingUsecase {
	return &PollingUsecase{
		OpcService: s,
		Repo:       r,
	}
}

// GetControlProgram - получение управляющей программы
func (u *PollingUsecase) GetControlProgram(req models.GetControlProgramRequest) (*models.ControlProgramInfoRequest, error) {
	if req.UUID == "" {
		return nil, fmt.Errorf("SessionID is empty")
	}
	id, err := uuid.Parse(req.UUID)
	if err != nil {
		return nil, err
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
func (u *PollingUsecase) StartPollingMachine(machineID uuid.UUID) *errors.AppError {
	err := u.OpcService.StartPollingForMachine(machineID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return errors.NewAppError(http.StatusNotFound, "machine not found", err, false)
		} else {
			return errors.NewAppError(http.StatusInternalServerError, "failed to start polling for machine", err, false)
		}
	}
	updateMap := map[string]interface{}{
		"status": connection_models.ConnectionStatusPolled,
	}

	_, err = u.Repo.UpdateCncMachine(machineID.String(), updateMap)
	if err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "failed to update machine record", err, false)
	}
	return nil
}

// StopPollingMachine останавливает сбор данных машины по UUID
func (u *PollingUsecase) StopPollingMachine(machineID uuid.UUID) *errors.AppError {
	err := u.OpcService.StopPollingForMachine(machineID)
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

	_, err = u.Repo.UpdateCncMachine(machineID.String(), updateMap)
	if err != nil {
		return errors.NewAppError(http.StatusInternalServerError, "failed to update machine record", err, false)
	}
	return nil
}

package interfaces

import (
	"opc_ua_service/internal/domain/entities"
)

type Repository interface {
	CncMachineRepository
}

type CncMachineRepository interface {
	CreateCncMachine(cnc *entities.CncMachine) (string, error)
	GetCncMachineBySIK(sik string) (entities.CncMachine, error)
	UpdateCncMachine(sik string, updateMap map[string]interface{}) (string, error)
	DeleteCncMachine(sik string) error
	GetAllCncMachines() ([]entities.CncMachine, error)
}

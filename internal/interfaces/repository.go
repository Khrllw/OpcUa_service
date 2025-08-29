package interfaces

import (
	"opc_ua_service/internal/domain/entities"
)

type Repository interface {
	CncMachineRepository
}

type CncMachineRepository interface {
	CreateCncMachine(cnc *entities.CncMachine) (string, error)
	GetCncMachineByUUID(id string) (entities.CncMachine, error)
	GetCncMachineByEndpointURL(endpoint string) (entities.CncMachine, error)
	UpdateCncMachine(id string, updateMap map[string]interface{}) (string, error)
	DeleteCncMachine(id string) error
	DeleteAllCncMachines() error
	GetAllCncMachines() ([]entities.CncMachine, error)
}

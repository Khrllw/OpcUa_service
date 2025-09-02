package interfaces

import (
	"opc_ua_service/internal/domain/entities"
)

type Repository interface {
	CncMachineRepository
	PasswordConnectionRepository
	AnonymousConnectionRepository
	CertificateConnectionRepository
}

type CncMachineRepository interface {
	CreateCncMachine(cnc entities.CncMachine) (string, error)
	GetCncMachineByUUID(id string) (entities.CncMachine, error)
	GetCncMachineByEndpointURL(endpoint string) (entities.CncMachine, error)
	UpdateCncMachine(id string, updateMap map[string]interface{}) (string, error)
	DeleteCncMachine(id string) error
	DeleteAllCncMachines() error
	GetAllCncMachines() ([]entities.CncMachine, error)
}

type PasswordConnectionRepository interface {
	CreatePasswordConnection(pc entities.PasswordConnection) (uint, error)
	GetPasswordConnectionByID(id uint) (entities.PasswordConnection, error)
	UpdatePasswordConnection(id uint, updateMap map[string]interface{}) (uint, error)
	DeletePasswordConnection(id uint) error
}

type AnonymousConnectionRepository interface {
	CreateAnonymousConnection(ac entities.AnonymousConnection) (uint, error)
	GetAnonymousConnectionByID(id uint) (entities.AnonymousConnection, error)
	UpdateAnonymousConnection(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteAnonymousConnection(id uint) error
}

type CertificateConnectionRepository interface {
	CreateCertificateConnection(cc entities.CertificateConnection) (uint, error)
	GetCertificateConnectionByID(id uint) (entities.CertificateConnection, error)
	UpdateCertificateConnection(id uint, updateMap map[string]interface{}) (uint, error)
	DeleteCertificateConnection(id uint) error
}

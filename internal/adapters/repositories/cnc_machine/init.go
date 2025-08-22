package cnc_machine

import (
	"gorm.io/gorm"
	"opc_ua_service/internal/interfaces"
)

type CncMachineRepositoryImpl struct {
	db *gorm.DB
}

func NewCncMachineRepository(db *gorm.DB) interfaces.CncMachineRepository {
	return &CncMachineRepositoryImpl{db: db}
}

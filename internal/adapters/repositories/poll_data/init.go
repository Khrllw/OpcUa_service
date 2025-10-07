package poll_data

import (
	"gorm.io/gorm"
	"opc_ua_service/internal/interfaces"
)

type PollDataRepositoryImpl struct {
	db *gorm.DB
}

func NewPollDataRepository(db *gorm.DB) interfaces.PollDataRepository {
	return &PollDataRepositoryImpl{db: db}
}

package password_connection

import (
	"gorm.io/gorm"
	"opc_ua_service/internal/interfaces"
)

type PasswordConnectionRepositoryImpl struct {
	db *gorm.DB
}

func NewPasswordConnectionRepository(db *gorm.DB) interfaces.PasswordConnectionRepository {
	return &PasswordConnectionRepositoryImpl{db: db}
}

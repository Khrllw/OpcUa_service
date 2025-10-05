package entities

import (
	connection_models "opc_ua_service/internal/domain/models/connection_types"
	"time"
)

// CncMachine представляет информацию о контроллере станка с ЧПУ для БД
type CncMachine struct {
	ID             uint64                                 `gorm:"primary_key;auto_increment"`
	UUID           string                                 `gorm:"not null;unique" json:"session_id"`
	EndpointURL    string                                 `gorm:"not null" json:"endpoint_url"`
	Model          string                                 `gorm:"not null" json:"model"`
	Manufacturer   string                                 `json:"manufacturer"`
	CreatedAt      time.Time                              `json:"created_at"`
	UpdatedAt      time.Time                              `json:"updated_at"`
	Status         connection_models.ConnectionStatusEnum `gorm:"not null" json:"status"`          // connected / polled
	PollTimeout    int                                    `json:"poll_timeout"`                    // Интервал опроса в сек
	ConnectionType connection_models.ConnectionTypeEnum   `gorm:"not null" json:"connection_type"` // "certificate", "anonymous", "password"

	CertificateConnectionID *uint
	AnonymousConnectionID   *uint
	PasswordConnectionID    *uint

	// Опциональные связи
	CertificateConnection *CertificateConnection `gorm:"foreignKey:CertificateConnectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	AnonymousConnection   *AnonymousConnection   `gorm:"foreignKey:AnonymousConnectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	PasswordConnection    *PasswordConnection    `gorm:"foreignKey:PasswordConnectionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
}

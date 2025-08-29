package entities

import "time"

type ConnectionStatus string

const (
	ConnectionStatusConnected    ConnectionStatus = "connected"
	ConnectionStatusDisconnected ConnectionStatus = "disconnected"
	ConnectionStatusPolled       ConnectionStatus = "polled"
)

// CncMachine представляет информацию о контроллере станка с ЧПУ для БД
type CncMachine struct {
	UUID         string           `gorm:"primaryKey;not null" json:"session_id"`
	EndpointURL  string           `gorm:"not null" json:"endpoint_url"`
	Model        string           `gorm:"not null" json:"model"`
	Manufacturer string           `json:"manufacturer"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	Status       ConnectionStatus `gorm:"not null" json:"status"` // connected / disconnected / polled
	Interval     int              `json:"interval"`               // Интервал опроса в сек
}

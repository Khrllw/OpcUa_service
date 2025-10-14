package entities

import (
	"gorm.io/datatypes"
	"time"
)

type PollData struct {
	ID        uint           `gorm:"primaryKey"`
	Timestamp time.Time      `gorm:"index"`
	Payload   datatypes.JSON `gorm:"type:jsonb"`
	Processed bool           `gorm:"default:false;index"`
}

package entities

import "time"

// CncMachine представляет информацию о контроллере станка с ЧПУ
type CncMachine struct {
	SIK       string    `gorm:"primaryKey;not null" json:"sik" example:"1234-5678-ABCD" rus:"SIK (идентификатор системы)"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	//Name string `gorm:"not null" json:"name" example:"Фрезерный станок №1" rus:"Название"`

	ControllerManufacturer string `gorm:"not null" json:"manufacturer" example:"DR. JOHANNES HEIDENHAIN GmbH" rus:"Производитель контроллера"`
	ControllerModel        string `gorm:"not null" json:"model" example:"TNC640" rus:"Модель контроллера"`
}

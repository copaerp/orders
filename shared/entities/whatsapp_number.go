package entities

import (
	"time"

	"github.com/google/uuid"
)

type WhatsappNumber struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	UnitID       uuid.UUID `gorm:"type:char(36);not null"`
	Number       string
	Description  string
	MetaNumberID string

	CreatedAt time.Time
	UpdatedAt time.Time

	Unit *Unit `gorm:"foreignKey:UnitID"`
}

func (w *WhatsappNumber) TableName() string {
	return "whatsapp_number"
}

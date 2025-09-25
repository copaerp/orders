package entities

import (
	"time"

	"github.com/google/uuid"
)

type WhatsappNumber struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	UnitID       uuid.UUID `gorm:"type:char(36);not null" json:"unit_id"`
	Number       string    `json:"number"`
	Description  string    `json:"description"`
	MetaNumberID string    `json:"meta_number_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Unit *Unit `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}

func (w *WhatsappNumber) TableName() string {
	return "whatsapp_number"
}

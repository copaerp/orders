package entities

import (
	"time"

	"github.com/google/uuid"
)

type WhatsappNumber struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	BusinessID uuid.UUID `gorm:"type:char(36);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Business *Business `gorm:"foreignKey:BusinessID"`
	Orders   []Order   `gorm:"foreignKey:WhatsappNumberID"`
}

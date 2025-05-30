package entities

import (
	"time"

	"github.com/google/uuid"
)

type WhatsappNumber struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	UnitID      uuid.UUID `gorm:"type:char(36);not null"`
	Number      string    `gorm:"type:varchar(20);not null"`
	Description string    `gorm:"type:varchar(255);default:null"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Unit *Unit `gorm:"foreignKey:UnitID"`
}

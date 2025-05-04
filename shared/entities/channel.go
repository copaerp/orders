package entities

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time

	Orders []Order
}

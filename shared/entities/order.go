package entities

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey"`
	CustomerID    uuid.UUID `gorm:"type:char(36);not null"`
	UnitID        uuid.UUID `gorm:"type:char(36);not null"`
	ChannelID     uuid.UUID `gorm:"type:char(36);not null"`
	Status        string
	Notes         string
	PaymentMethod string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FinishedAt    *time.Time
	CanceledAt    *time.Time

	Customer       Customer `gorm:"foreignKey:CustomerID"`
	Unit           Unit     `gorm:"foreignKey:UnitID"`
	Channel        Channel  `gorm:"foreignKey:ChannelID"`
	ProductsOrders []ProductOrder
}

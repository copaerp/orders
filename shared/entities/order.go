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

	Customer       Customer
	Unit           Unit
	Channel        Channel
	ProductsOrders []ProductOrder
}

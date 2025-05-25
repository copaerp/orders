package entities

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	BusinessID  uuid.UUID `gorm:"type:char(36);not null"`
	Name        string
	Description string
	BRLPrice    float64
	Category    string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Business        *Business `gorm:"foreignKey:BusinessID"`
	ProductsInUnits []ProductInUnit
	ProductsOrders  []ProductOrder
}

package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProductOrder struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	ProductID uuid.UUID `gorm:"type:char(36);not null"`
	OrderID   uuid.UUID `gorm:"type:char(36);not null"`
	Amount    int
	CreatedAt time.Time
	UpdatedAt time.Time

	Product Product `gorm:"foreignKey:ProductID"`
	Order   Order   `gorm:"foreignKey:OrderID"`
}

package entities

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID         uuid.UUID `gorm:"type:char(36);primaryKey"`
	BusinessID uuid.UUID `gorm:"type:char(36);not null"`
	FullName   string
	Phone      string
	Email      string
	Document   string
	BirthDate  time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time

	Business Business `gorm:"foreignKey:BusinessID"`
	Orders   []Order
}

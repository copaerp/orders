package entities

import (
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey"`
	BusinessID   uuid.UUID `gorm:"type:char(36);not null"`
	Name         string
	Phone        string
	PostalCode   string
	StreetName   string
	StreetNumber string
	City         string
	State        string
	Country      string
	Neighborhood string
	Complement   string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Business        *Business `gorm:"foreignKey:BusinessID"`
	ProductsInUnits []ProductInUnit
	Orders          []Order
	WhatsappNumbers []WhatsappNumber
}

func (u *Unit) TableName() string {
	return "unit"
}

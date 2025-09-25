package entities

import (
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	ID           uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessID   uuid.UUID `gorm:"type:char(36);not null" json:"business_id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	PostalCode   string    `json:"postal_code"`
	StreetName   string    `json:"street_name"`
	StreetNumber string    `json:"street_number"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	Country      string    `json:"country"`
	Neighborhood string    `json:"neighborhood"`
	Complement   string    `json:"complement"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Business        *Business       `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	WhatsappNumber  *WhatsappNumber `gorm:"foreignKey:UnitID" json:"whatsapp_number,omitempty"`
	ProductsInUnits []ProductInUnit `json:"products_in_units,omitempty"`
	Orders          []Order         `json:"orders,omitempty"`
}

func (u *Unit) TableName() string {
	return "unit"
}

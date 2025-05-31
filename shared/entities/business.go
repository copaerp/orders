package entities

import (
	"time"

	"github.com/google/uuid"
)

type Business struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name      string
	LegalName string
	CNPJ      string
	Email     string
	Phone     string
	LogoURL   string
	Industry  string
	CreatedAt time.Time
	UpdatedAt time.Time

	Units     []Unit
	Products  []Product
	Customers []Customer
}

func (b *Business) TableName() string {
	return "business"
}

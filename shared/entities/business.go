package entities

import (
	"time"

	"github.com/google/uuid"
)

type Business struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `json:"name"`
	LegalName string    `json:"legal_name"`
	CNPJ      string    `json:"cnpj"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	LogoURL   string    `json:"logo_url"`
	Industry  string    `json:"industry"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Units     []Unit     `json:"units,omitempty"`
	Products  []Product  `json:"products,omitempty"`
	Customers []Customer `json:"customers,omitempty"`
}

func (b *Business) TableName() string {
	return "business"
}

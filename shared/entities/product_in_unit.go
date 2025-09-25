package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProductInUnit struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID uuid.UUID `gorm:"type:char(36);not null" json:"product_id"`
	UnitID    uuid.UUID `gorm:"type:char(36);not null" json:"unit_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Unit    *Unit    `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
}

func (p *ProductInUnit) TableName() string {
	return "product_in_unit"
}

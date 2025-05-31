package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProductInUnit struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	ProductID uuid.UUID `gorm:"type:char(36);not null"`
	UnitID    uuid.UUID `gorm:"type:char(36);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Product *Product `gorm:"foreignKey:ProductID"`
	Unit    *Unit    `gorm:"foreignKey:UnitID"`
}

func (p *ProductInUnit) TableName() string {
	return "product_in_unit"
}

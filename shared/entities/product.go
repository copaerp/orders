package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:char(36);primaryKey"`
	BusinessID  uuid.UUID       `gorm:"type:char(36);not null;index"`
	Name        string          `gorm:"type:varchar(255);not null"`
	Description string          `gorm:"type:text"`
	BRLPrice    decimal.Decimal `gorm:"type:decimal(19,4);not null"`
	Category    string          `gorm:"type:varchar(100)"`
	ImageURL    string          `gorm:"type:varchar(255)"`
	CreatedAt   time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP"`

	Business        *Business       `gorm:"foreignKey:BusinessID;references:ID"`
	ProductsInUnits []ProductInUnit `gorm:"foreignKey:ProductID"`
	ProductsOrders  []ProductOrder  `gorm:"foreignKey:ProductID"`
}

func (p *Product) TableName() string {
	return "product"
}

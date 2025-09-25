package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	ID          uuid.UUID       `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessID  uuid.UUID       `gorm:"type:char(36);not null;index" json:"business_id"`
	Name        string          `gorm:"type:varchar(255);not null" json:"name"`
	Description string          `gorm:"type:text" json:"description"`
	BRLPrice    decimal.Decimal `gorm:"type:decimal(19,4);not null" json:"brl_price"`
	Category    string          `gorm:"type:varchar(100)" json:"category"`
	ImageURL    string          `gorm:"type:varchar(255)" json:"image_url"`
	CreatedAt   time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"type:timestamp;default:CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP" json:"updated_at"`

	Business        *Business       `gorm:"foreignKey:BusinessID;references:ID" json:"business,omitempty"`
	ProductsInUnits []ProductInUnit `gorm:"foreignKey:ProductID" json:"products_in_units,omitempty"`
	ProductsOrders  []ProductOrder  `gorm:"foreignKey:ProductID" json:"products_orders,omitempty"`
}

func (p *Product) TableName() string {
	return "product"
}

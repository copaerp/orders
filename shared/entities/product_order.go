package entities

import (
	"time"

	"github.com/google/uuid"
)

type ProductOrder struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	ProductID uuid.UUID `gorm:"type:char(36);not null" json:"product_id"`
	OrderID   uuid.UUID `gorm:"type:char(36);not null" json:"order_id"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Order   *Order   `gorm:"foreignKey:OrderID" json:"order,omitempty"`
}

func (p *ProductOrder) TableName() string {
	return "product_order"
}

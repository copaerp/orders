package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Order struct {
	ID                 uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	DisplayID          string         `gorm:"type:varchar(12);not null" json:"display_id"`
	CustomerID         uuid.UUID      `gorm:"type:char(36);not null" json:"customer_id"`
	UnitID             uuid.UUID      `gorm:"type:char(36);not null" json:"unit_id"`
	ChannelID          uuid.UUID      `gorm:"type:char(36);not null" json:"channel_id"`
	Status             string         `json:"status"`
	PostCheckoutStatus string         `json:"post_checkout_status"`
	Notes              *string        `json:"notes,omitempty"`
	PaymentMethod      *string        `json:"payment_method,omitempty"`
	UsedMenu           []byte         `gorm:"type:blob;default:NULL" json:"used_menu,omitempty"`
	CurrentCart        datatypes.JSON `gorm:"type:json;default:NULL" json:"current_cart,omitempty"`
	LastMessageAt      time.Time      `json:"last_message_at"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	FinishedAt         *time.Time     `json:"finished_at,omitempty"`
	CanceledAt         *time.Time     `json:"canceled_at,omitempty"`

	Customer       *Customer      `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Unit           *Unit          `gorm:"foreignKey:UnitID" json:"unit,omitempty"`
	Channel        *Channel       `gorm:"foreignKey:ChannelID" json:"channel,omitempty"`
	ProductsOrders []ProductOrder `json:"products_orders,omitempty"`
}

func (o *Order) TableName() string {
	return "order"
}

func (o *Order) GetMenuAsMapArr() ([]map[string]string, error) {
	var menu []map[string]string
	if err := json.Unmarshal(o.UsedMenu, &menu); err != nil {
		return nil, err
	}
	return menu, nil
}

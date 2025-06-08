package entities

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey"`
	CustomerID    uuid.UUID `gorm:"type:char(36);not null"`
	UnitID        uuid.UUID `gorm:"type:char(36);not null"`
	ChannelID     uuid.UUID `gorm:"type:char(36);not null"`
	Status        string
	Notes         *string
	PaymentMethod *string
	UsedMenu      []byte `gorm:"type:blob;default:NULL"`
	LastMessageAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
	FinishedAt    *time.Time
	CanceledAt    *time.Time

	Customer       *Customer `gorm:"foreignKey:CustomerID"`
	Unit           *Unit     `gorm:"foreignKey:UnitID"`
	Channel        *Channel  `gorm:"foreignKey:ChannelID"`
	ProductsOrders []ProductOrder
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

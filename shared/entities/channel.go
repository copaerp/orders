package entities

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Orders []Order `json:"orders,omitempty"`
}

func (c *Channel) TableName() string {
	return "channel"
}

package entities

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID            uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	BusinessID    uuid.UUID `gorm:"type:char(36);not null" json:"business_id"`
	FullName      string    `json:"full_name"`
	Phone         string    `json:"phone"`
	InstagramUser string    `json:"instagram_user"`
	Email         string    `json:"email"`
	Document      string    `json:"document"`
	BirthDate     time.Time `json:"birth_date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Business *Business `gorm:"foreignKey:BusinessID" json:"business,omitempty"`
	Orders   []Order   `json:"orders,omitempty"`
}

func (c *Customer) TableName() string {
	return "customer"
}

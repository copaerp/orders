package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Business struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey"`
	Name      string         `gorm:"type:varchar(255);not null"`
	LegalName string         `gorm:"type:varchar(255)"`
	CNPJ      string         `gorm:"type:char(14);uniqueIndex"`
	Email     string         `gorm:"type:varchar(255)"`
	Phone     string         `gorm:"type:varchar(20)"`
	LogoURL   string         `gorm:"type:text"`
	Industry  string         `gorm:"type:varchar(100)"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

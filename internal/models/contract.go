package models

import (
	"time"

	"gorm.io/gorm"
)

type Contract struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	TenderUUID string `gorm:"type:uuid"`
	VendorUUID string `gorm:"type:uuid"`
	StartDate  time.Time
	EndDate    time.Time
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type PurchaseOrder struct {
	ID           string `gorm:"type:uuid;primaryKey"`
	ContractUUID string `gorm:"type:uuid"`
	TotalAmount  float64
	Status       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

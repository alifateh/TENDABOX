package models

import (
	"time"

	"gorm.io/gorm"
)

type Vendor struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	UserUUID   string `gorm:"type:uuid"`
	Name       string
	IsPublic   bool
	IsVerified bool
	Score      float64
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type VendorSubscription struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	VendorUUID string `gorm:"type:uuid"`
	Status     string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

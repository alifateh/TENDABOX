package models

import (
	"time"

	"gorm.io/gorm"
)

type PurchaseRequest struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	RequestedBy string `gorm:"type:uuid"`
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type Tender struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	PRUUID    string `gorm:"type:uuid"`
	Title     string
	Type      string // public / premium
	Status    string // draft/published/closed
	Deadline  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

type TenderItem struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	TenderUUID string `gorm:"type:uuid"`
	Name       string
	Quantity   int
	UnitPrice  float64
}

package models

import (
	"time"

	"gorm.io/gorm"
)

type Bid struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	TenderUUID string `gorm:"type:uuid"`
	VendorUUID string `gorm:"type:uuid"`
	Status     string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt
}

type BidEvaluation struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	BidUUID   string `gorm:"type:uuid"`
	ScoreTech float64
	ScoreFin  float64
	AIFlag    bool
	CreatedAt time.Time
}

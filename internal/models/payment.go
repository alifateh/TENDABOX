package models

import "time"

type Invoice struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	POUUID    string `gorm:"type:uuid"`
	Amount    float64
	Status    string
	CreatedAt time.Time
}

type Payment struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	InvoiceUUID string `gorm:"type:uuid"`
	Method      string // MPesa / Visa
	Status      string
	PaidAt      time.Time
}

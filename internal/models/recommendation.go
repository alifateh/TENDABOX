package models

import "time"

type AIScore struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	Entity     string // tender/vendor/bid
	EntityUUID string
	Score      float64
	Reason     string
	CreatedAt  time.Time
}

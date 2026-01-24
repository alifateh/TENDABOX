package models

import "time"

type AuditLog struct {
	ID         string `gorm:"type:uuid;primaryKey"`
	UserUUID   string `gorm:"type:uuid"`
	Action     string
	Entity     string
	EntityUUID string
	CreatedAt  time.Time
}

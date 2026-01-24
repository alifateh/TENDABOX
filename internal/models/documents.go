package models

import "time"

type Document struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	FilePath  string
	Version   int
	Hash      string
	CreatedAt time.Time
}

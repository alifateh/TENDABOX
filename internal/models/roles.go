package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Roles struct {
	ID        string         `gorm:"type:uuid;primaryKey;" json:"id"`
	Name      string         `gorm:"column:role_name;uniqueIndex;" json:"role_name"`
	Level     uint           `gorm:"column:role_level" json:"role_level"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (r *Roles) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = uuid.NewString()
	return
}

type Permission struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	Code      string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

func (p *Permission) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.NewString()
	return
}

type RolePermission struct {
	RoleUUID       string `gorm:"type:uuid"`
	PermissionUUID string `gorm:"type:uuid"`
}

func (rp *RolePermission) BeforeCreate(tx *gorm.DB) (err error) {
	rp.RoleUUID = uuid.NewString()
	return
}

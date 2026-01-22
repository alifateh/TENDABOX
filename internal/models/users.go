package models

import (
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;primarykey;" json:"id"`
	FirstName string         `gorm:"column:first_name" json:"first_name"`
	LastName  string         `gorm:"column:last_name" json:"last_name"`
	Email     string         `gorm:"unique;index;not null;" json:"email"`
	Password  string         `gorm:"not null;" json:"-"`
	IsActive  bool           `gorm:"column:is_active" json:"is_active"`
	Role      Roles          `gorm:"foreignKey:RoleUUID"`
	RoleUUID  string         `gorm:"type:uuid"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Warn("This is a Fatal-Error in genrating Password Encryption", "Error", err)
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(providedPassword))
}

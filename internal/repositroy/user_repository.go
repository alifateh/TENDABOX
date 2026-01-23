package repositroy

import (
	"Tendabox/internal/models"

	"gorm.io/gorm"
)

// UserRepository تعریف رفتارها
type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	// می‌توانید متدهای دیگر مثل Create یا Update را اینجا اضافه کنید
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository سازنده ریپازیتوری
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// GetByEmail پیاده‌سازی کوئری دیتابیس
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

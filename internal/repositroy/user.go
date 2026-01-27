package repositroy

import (
	"Tendabox/internal/models"
	"fmt"

	"gorm.io/gorm"
)

// UserRepository تعریف رفتارها
type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	CreateUser(user *models.User) error
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
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	// ۱. بررسی اینکه آیا ایمیل از قبل وجود دارد یا خیر
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)

	if count > 0 {
		// بازگرداندن یک خطای سفارشی که در Handler قابل شناسایی باشد
		return fmt.Errorf("user_already_exists")
	}

	// ۲. تلاش برای ایجاد کاربر جدید
	// نکته: &user را به user تغییر دادم چون خودش پوینتر است
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

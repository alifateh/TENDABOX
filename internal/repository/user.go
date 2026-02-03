package repositroy

import (
	"Tendabox/internal/models"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

// UserRepository
type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	GetAllUsers(page int, pageSize int) ([]models.User, int64, error)
	CreateUser(user *models.User) error
	UpdateUserRole(UserID string, RoleUUID string) (err error)
	UpdateUserStatus(UserID string, NewStatus bool) (err error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository سازنده ریپازیتوری
func NewUserRepository(DB *gorm.DB) UserRepository {
	return &userRepository{db: DB}
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ?", email).Where("is_active = ?", true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) CreateUser(user *models.User) error {
	var count int64
	r.db.Model(&models.User{}).Where("email = ?", user.Email).Count(&count)

	if count > 0 {
		return fmt.Errorf("user_already_exists")
	}

	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) UpdateUserRole(UserID string, RoleUUID string) (err error) {

	result := r.db.Model(&models.User{}).
		Where("id = ?", UserID).
		Update("role_uuid", RoleUUID)

	if result.Error != nil {
		slog.Error("Failed to update user role", "userID", UserID, "error", result.Error)
		return result.Error
	}
	if result.RowsAffected == 0 {
		slog.Warn("No user found to update role", "userID", UserID)
		return fmt.Errorf("user not found")
	}
	return nil

}

func (r *userRepository) GetAllUsers(page int, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	r.db.Model(&models.User{}).Count(&total)

	offset := (page - 1) * pageSize

	err := r.db.Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) UpdateUserStatus(UserID string, NewStatus bool) error {
	// Correct GORM syntax: Update("column_name", value)
	result := r.db.Model(&models.User{}).
		Where("id = ?", UserID).
		Update("is_active", NewStatus).Error // Removed "=" and "?" from the string

	if result != nil {
		slog.Error("Update User status failed", "Error", result)
		return result
	}

	slog.Info("Update Success", "UserID", UserID, "NewStatus", NewStatus)
	return nil
}

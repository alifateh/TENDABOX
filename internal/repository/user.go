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
	GetAllUsers() ([]models.User, error)
	CreateUser(user *models.User) error
	UpdateUserRole(UserID string, RoleUUID string) (err error)
}

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository سازنده ریپازیتوری
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
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

func (r *userRepository) GetAllUsers() ([]models.User, error) {
	var Users []models.User
	err := r.db.Find(&Users).Error
	if err != nil {
		slog.Error("Fatal Error select all users")
		return nil, err
	}
	return Users, nil
}

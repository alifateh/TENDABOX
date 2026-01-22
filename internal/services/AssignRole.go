package services

import (
	"Tendabox/internal/models"

	"gorm.io/gorm"
)

func AssignRole(tx *gorm.DB, userID, roleID string) error {
	return tx.Create(&models.RolePermission{
		RoleUUID:       userID,
		PermissionUUID: roleID,
	}).Error
}

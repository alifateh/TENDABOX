package repositroy

import (
	"Tendabox/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

func AllRoles(db *gorm.DB) ([]models.Roles, error) {
	var All_Roles []models.Roles

	result := db.Find(&All_Roles)

	if result.Error != nil {
		slog.Warn("Error in DB Response", "Error", result.Error)
		return nil, result.Error
	} else {
		return All_Roles, nil
	}

}

package repositroy

import (
	"Tendabox/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

func GenrateMenu(RoleName string, db *gorm.DB) (Main_Menu []models.MainMenu, err error) {
	var role models.Roles
	var menus []models.MainMenu
	result := db.Where("role_name=?", RoleName).Find(&role)
	if result.Error != nil {
		slog.Warn("User Role not Found!", "Fatal Error", result.Error)
		return nil, result.Error
	}

	Menu := db.Where("min_level <= ?", role.Level).Order("min_level ASC").Find(&menus)
	if Menu.Error != nil {
		slog.Warn("No Menu Found!", "Fatal Error", Menu.Error)
		return nil, Menu.Error
	} else {
		return menus, nil
	}
}

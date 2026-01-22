package main

import (
	"Tendabox/internal/logger"
	"Tendabox/internal/models"
	"Tendabox/pkg/database"
	"fmt"
	"log/slog"

	"gorm.io/gorm"
)

func main() {
	logger.SetupLogger()
	database.Connect()
	flag := 0

	err := database.DB.AutoMigrate(&models.User{})
	if err != nil {
		slog.Warn("Failed to migrate Table [User]", "Error", err)
		flag++
	} else {
		slog.Info("Migration of Table [User] is Done")
	}

	err = database.DB.AutoMigrate(&models.Permission{})

	if err != nil {
		slog.Warn("Failed to migrate Table [Permission]", "Error", err)
		flag++
	} else {
		slog.Info("Migration of Table [Permission] is Done")
	}

	err = database.DB.AutoMigrate(&models.Roles{})

	if err != nil {
		slog.Warn("Failed to migrate Table [Roles]", "Error", err)
		flag++
	} else {
		slog.Info("Migration of Table [Roles] is Done")
	}

	err = database.DB.AutoMigrate(&models.RolePermission{})

	if err != nil {
		slog.Warn("Failed to migrate Table [RolePermission]", "Error", err)
		flag++
	} else {
		slog.Info("Migration of Table [RolePermission] is Done")
	}

	if flag == 0 {
		err = RunSeed(database.DB)
		if err != nil {
			slog.Warn("Failed to Seed", err)
		} else {
			slog.Info("All Seeding is Done")
		}

	}

}

func SeedRoles(db *gorm.DB) (map[string]string, error) {
	roles := []models.Roles{
		{Name: "super_admin", Level: 2000},
		{Name: "admin", Level: 1000},
		{Name: "procurement_manager", Level: 800},
		{Name: "finance", Level: 700},
		{Name: "auditor", Level: 600},
		{Name: "vendor_premium", Level: 500},
		{Name: "buyer", Level: 400},
		{Name: "vendor_public", Level: 100},
		{Name: "viewer", Level: 50},
	}

	roleMap := make(map[string]string)

	for _, r := range roles {
		var role models.Roles
		err := db.Where("role_name = ?", r.Name).First(&role).Error

		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}

			role = r
			if err := db.Create(&role).Error; err != nil {
				return nil, err
			}
		}

		roleMap[role.Name] = role.ID
	}

	return roleMap, nil
}

func SeedPermissions(db *gorm.DB) (map[string]string, error) {
	permissions := []models.Permission{
		{Code: "user.manage"},
		{Code: "role.manage"},

		{Code: "tender.create"},
		{Code: "tender.publish"},
		{Code: "tender.close"},
		{Code: "tender.view"},

		{Code: "bid.submit"},
		{Code: "bid.view"},
		{Code: "bid.evaluate"},

		{Code: "contract.create"},
		{Code: "contract.sign"},
		{Code: "contract.view"},

		{Code: "po.create"},
		{Code: "po.view"},

		{Code: "invoice.create"},
		{Code: "invoice.view"},
		{Code: "payment.process"},

		{Code: "document.upload"},
		{Code: "document.view"},

		{Code: "audit.view"},
	}

	permMap := make(map[string]string)

	for _, p := range permissions {
		var perm models.Permission

		err := db.Where("code = ?", p.Code).First(&perm).Error
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				return nil, err
			}
			perm = p
			if err := db.Create(&perm).Error; err != nil {
				return nil, err
			}
		}

		permMap[perm.Code] = perm.ID
	}

	return permMap, nil
}

func SeedRolePermissions(
	db *gorm.DB,
	roleMap map[string]string,
	permMap map[string]string,
) error {

	for roleName, permCodes := range rolePermissions {
		roleID, ok := roleMap[roleName]
		if !ok {
			return fmt.Errorf("role not found: %s", roleName)
		}

		for _, code := range permCodes {
			permID, ok := permMap[code]
			if !ok {
				return fmt.Errorf("permission not found: %s", code)
			}

			rp := models.RolePermission{
				RoleUUID:       roleID,
				PermissionUUID: permID,
			}

			err := db.FirstOrCreate(
				&rp,
				models.RolePermission{
					RoleUUID:       roleID,
					PermissionUUID: permID,
				},
			).Error

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func SeedUsers(db *gorm.DB, roleMap map[string]string) error {
	users := []models.User{
		{
			FirstName: "Admin",
			LastName:  "User",
			Email:     "admin@tendabox.org",
			Password:  "Admin@123",
			RoleUUID:  roleMap["admin"],
		},
		{
			FirstName: "Public",
			LastName:  "Vendor",
			Email:     "vendor_public@tendabox.org",
			Password:  "Vendor@123",
			RoleUUID:  roleMap["vendor_public"],
		},
		{
			FirstName: "Premium",
			LastName:  "Vendor",
			Email:     "vendor_premium@tendabox.org",
			Password:  "Vendor@123",
			RoleUUID:  roleMap["vendor_premium"],
		},
	}

	for _, u := range users {
		db.FirstOrCreate(&u, models.User{Email: u.Email})
	}

	return nil
}

func RunSeed(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {

		roleMap, err := SeedRoles(tx)
		if err != nil {
			return err
		}

		permMap, err := SeedPermissions(tx)
		if err != nil {
			return err
		}

		if err := SeedRolePermissions(tx, roleMap, permMap); err != nil {
			return err
		}

		if err := SeedUsers(tx, roleMap); err != nil {
			return err
		}

		return nil
	})
}

var rolePermissions = map[string][]string{
	"super_admin": {
		"user.manage", "role.manage",
		"tender.create", "tender.publish", "tender.close", "tender.view",
		"bid.submit", "bid.view", "bid.evaluate",
		"contract.create", "contract.sign", "contract.view",
		"po.create", "po.view",
		"invoice.create", "invoice.view", "payment.process",
		"document.upload", "document.view",
		"audit.view",
	},
	"admin": {
		"user.manage",
		"tender.create", "tender.publish", "tender.close", "tender.view",
		"bid.view", "bid.evaluate",
		"contract.create", "contract.sign", "contract.view",
		"po.create", "po.view",
		"invoice.view", "payment.process",
		"audit.view",
	},
	"procurement_manager": {
		"tender.create", "tender.publish", "tender.close", "tender.view",
		"bid.view", "bid.evaluate",
		"contract.create", "contract.view",
	},
	"finance": {
		"invoice.create", "invoice.view",
		"payment.process",
		"audit.view",
	},
	"auditor": {
		"audit.view",
		"tender.view",
		"bid.view",
		"contract.view",
		"invoice.view",
	},
	"vendor_premium": {
		"tender.view",
		"bid.submit",
		"document.upload",
	},
	"vendor_public": {
		"tender.view",
	},
	"buyer": {
		"tender.view",
		"bid.view",
	},
	"viewer": {
		"tender.view",
	},
}

//////////////////////////////////////
//////////////////RBAC////////////////////
//////////////////////////////////////

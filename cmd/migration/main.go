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

	err_menu := database.DB.AutoMigrate(&models.MainMenu{})
	if err_menu != nil {
		slog.Warn("Failed to migrate Table [Main Menu]", "Error", err_menu)
		return
	} else {
		slog.Info("Migration of Table [Main Menu] is Done")
		err_seed := SeedMenus(database.DB)
		if err_seed != nil {
			slog.Warn("Failed to Seed Table [Main Menu]", "Error", err_seed)
		} else {
			slog.Info("Seeding of Table [Main Menu] is Done")
		}
	}

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
			slog.Warn("Failed to Seed", "Fatal_Error:", err)
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

		err := db.FirstOrCreate(
			&role,
			models.Roles{Name: r.Name},
		).Error
		if err != nil {
			return nil, err
		}

		// اگر تازه ساخته شده، level را ست کن
		db.Model(&role).Update("role_level", r.Level)

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
//////////////////Main_Menu////////////////////
//////////////////////////////////////

func SeedMenus(db *gorm.DB) error {
	menus := []models.MainMenu{
		// --- Standalone Items ---
		{ItemName: "Dashboard", URLPath: "/dashboard", MinLevel: 50, Icon: "bi bi-speedometer", ParentName: ""},

		// --- System Admin Group (Level 2000) ---
		{ItemName: "User Security & Roles", URLPath: "/admin/security", MinLevel: 2000, Icon: "bi bi-shield-lock", ParentName: "System Control"},
		{ItemName: "System Health & Logs", URLPath: "/admin/system/health", MinLevel: 2000, Icon: "bi bi-terminal", ParentName: "System Control"},
		{ItemName: "Global Settings", URLPath: "/admin/settings", MinLevel: 2000, Icon: "bi bi-gear", ParentName: "System Control"},

		// --- Admin Management Group (Level 1000) ---
		{ItemName: "User Management", URLPath: "/admin/users", MinLevel: 1000, Icon: "bi bi-people", ParentName: "Administration"},
		{ItemName: "Vendor Approval", URLPath: "/admin/vendors/approve", MinLevel: 1000, Icon: "bi bi-person-check", ParentName: "Administration"},
		{ItemName: "Executive Reports", URLPath: "/admin/reports", MinLevel: 1000, Icon: "bi bi-graph-up-arrow", ParentName: "Administration"},

		// --- Procurement Group (Level 800) ---
		{ItemName: "Tender Management", URLPath: "/procurement/tenders", MinLevel: 800, Icon: "bi bi-pennant", ParentName: "Procurement"},
		{ItemName: "Technical Evaluation", URLPath: "/procurement/evaluations", MinLevel: 800, Icon: "bi bi-clipboard-check", ParentName: "Procurement"},
		{ItemName: "AI Scoring Engine", URLPath: "/procurement/ai-score", MinLevel: 800, Icon: "bi bi-cpu", ParentName: "Procurement"},

		// --- Finance Group (Level 700) ---
		{ItemName: "Invoices & Payments", URLPath: "/finance/invoices", MinLevel: 700, Icon: "bi bi-receipt", ParentName: "Finance"},
		{ItemName: "Budget Control", URLPath: "/finance/budget", MinLevel: 700, Icon: "bi bi-wallet2", ParentName: "Finance"},
		{ItemName: "Financial Audit Trail", URLPath: "/finance/audit-trail", MinLevel: 700, Icon: "bi bi-currency-dollar", ParentName: "Finance"},

		// --- Auditor Group (Level 600) ---
		{ItemName: "Audit Dashboard", URLPath: "/audit/dashboard", MinLevel: 600, Icon: "bi bi-shredder", ParentName: "Audit"},
		{ItemName: "All Activity Logs", URLPath: "/audit/logs", MinLevel: 600, Icon: "bi bi-journal-text", ParentName: "Audit"},

		// --- Vendor Portal Group (Level 500 & 100) ---
		{ItemName: "Premium Tenders", URLPath: "/vendor/tenders-premium", MinLevel: 500, Icon: "bi bi-star-fill", ParentName: "Vendor Portal"},
		{ItemName: "My Proposals", URLPath: "/vendor/proposals", MinLevel: 500, Icon: "bi bi-send", ParentName: "Vendor Portal"},
		{ItemName: "Performance Analytics", URLPath: "/vendor/stats", MinLevel: 500, Icon: "bi bi-bar-chart-line", ParentName: "Vendor Portal"},
		{ItemName: "Public Tenders", URLPath: "/public/tenders", MinLevel: 100, Icon: "bi bi-megaphone", ParentName: "Vendor Portal"},
		{ItemName: "Company Profile", URLPath: "/vendor/profile", MinLevel: 100, Icon: "bi bi-building", ParentName: "Vendor Portal"},

		// --- Internal Buyer Group (Level 400) ---
		{ItemName: "My Purchase Requests", URLPath: "/buyer/requests", MinLevel: 400, Icon: "bi bi-cart", ParentName: "Purchasing"},
		{ItemName: "Create New PR", URLPath: "/buyer/requests/new", MinLevel: 400, Icon: "bi bi-cart-plus", ParentName: "Purchasing"},

		// --- Public/Viewer Group (Level 50) ---
		{ItemName: "Public Directory", URLPath: "/directory", MinLevel: 50, Icon: "bi bi-search", ParentName: "Public Info"},
		{ItemName: "General Statistics", URLPath: "/stats/public", MinLevel: 50, Icon: "bi bi-pie-chart", ParentName: "Public Info"},
	}

	for _, m := range menus {
		// برای جلوگیری از تکرار، از URLPath به عنوان کلید یکتا استفاده می‌کنیم
		err := db.Where(models.MainMenu{URLPath: m.URLPath}).
			Assign(models.MainMenu{
				ItemName:   m.ItemName,
				MinLevel:   m.MinLevel,
				Icon:       m.Icon,
				ParentName: m.ParentName,
			}).
			FirstOrCreate(&m).Error
		if err != nil {
			return err
		}
	}
	return nil
}

package routes

import (
	"net/http"

	"TENDABOX/internal/migration"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func registerMigrationRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/migration/setup", func(c *gin.Context) {

		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		logs, err := migration.EnsureTablesExistWithLogs(sqlDB)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "failed",
				"logs":   logs,
				"error":  err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"logs":   logs,
		})
	})
}

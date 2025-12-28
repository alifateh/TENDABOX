package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB) {

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "UP",
			"database": "connected",
			"time":     time.Now().Format(time.RFC3339),
		})
	})

	registerMigrationRoutes(r, db)
}

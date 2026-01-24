package routes

import (
	"Tendabox/internal/handlers"
	middleware "Tendabox/internal/middelwars"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("static", "static")

	// --- مسیرهای Frontend (رندر کردن صفحات) ---
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/dashboard", func(c *gin.Context) {
		// این صفحه در سمت کلاینت با JS چک می‌کند که آیا کوکی معتبر است یا خیر
		c.HTML(http.StatusOK, "NewDashboard.html", nil)
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/login", handlers.Login)
		//protected routes
		userGroup := v1.Group("/user")
		userGroup.Use(middleware.JWTAuth())
		{
			userGroup.GET("/dashboard-data", func(c *gin.Context) {
				userID, _ := c.Get("userID")
				userLevel, _ := c.Get("userLevel")
				slog.Info("Dashboard data requested", "user_id", userID, "User_Level", userLevel, "Client IP =", c.ClientIP())
				c.JSON(http.StatusOK, gin.H{
					"message": "خوش آمدید، اطلاعات با موفقیت بارگذاری شد",
					"status":  "authenticated",
					"level":   userLevel,
				})
			})

			// مسیرهای مخصوص ادمین و سوپر یوزر (RBAC)
			adminGroup := userGroup.Group("/admin")
			adminGroup.Use(middleware.AuthorizeRole("Admin", "Supper User"))
			{
				adminGroup.GET("/stats", func(c *gin.Context) {
					userID, _ := c.Get("userID")
					slog.Info("Admin stats accessed", "user_id", userID, "Client IP =", c.ClientIP())
					c.JSON(http.StatusOK, gin.H{
						"data": "آمار محرمانه سیستم برای مدیران",
					})
				})
			}
		}
	}

	return r
}

package routes

import (
	"Tendabox/internal/handlers"
	middleware "Tendabox/internal/middelwars"
	"Tendabox/pkg/database"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// --- مسیرهای Frontend (رندر کردن صفحات) ---
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})
	r.GET("/sample", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sample.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", nil)
	})

	r.NoRoute(func(c *gin.Context) {
		slog.Warn("URL Not Found", "Error 404", c.Request.RequestURI)
		if strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Endpoint Not Found"})
			return
		}
		c.HTML(http.StatusNotFound, "404.html", gin.H{"Error": "Endpoint Not Found"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.GET("/roles", handlers.GetAllRoles)
		userHandler := handlers.NewUserHandler(database.DB)

		v1.POST("/register", middleware.RegisterValidator, userHandler.RegisterUser)

		v1.POST("/login", handlers.Login)
		//protected routes
		userGroup := v1.Group("/user")
		userGroup.Use(middleware.JWTAuth())
		{
			userGroup.GET("/Accesslevel", handlers.GetAccessLevel)
			userGroup.GET("/MyMenu", handlers.GenerateMenu)

			// مسیرهای مخصوص ادمین و سوپر یوزر (RBAC)
			adminGroup := userGroup.Group("/admin")
			adminGroup.Use(middleware.AuthorizeRole("Admin", "super_admin"))
			{
				adminGroup.GET("/security", func(c *gin.Context) {
					userID, _ := c.Get("userID")
					slog.Info("Admin stats accessed", "user_id", userID, "Client IP =", c.ClientIP())
					c.HTML(200, "admin_security.html", gin.H{
						"data": "آمار محرمانه سیستم برای مدیران",
					})
				})
			}
		}

	}

	return r
}

package routes

import (
	"Tendabox/internal/handlers"
	middleware "Tendabox/internal/middelwars"
	repositroy "Tendabox/internal/repository"
	"Tendabox/pkg/database"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("./static", "./static")

	// --- مسیرهای Frontend (رندر کردن صفحات) ---
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "TendaBOX",
		})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard",
		})
	})
	r.GET("/sample", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sample.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Registeration",
		})
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

			//Admin Routes
			adminGroup := userGroup.Group("/admin")
			adminGroup.Use(middleware.AuthorizeRole("Admin", "super_admin"))
			{
				adminGroup.GET("/security", func(c *gin.Context) {
					c.HTML(200, "admin_security.html", gin.H{
						"title": "Manage User's Roles",
					})

				})
				userRepo := repositroy.NewUserRepository(database.DB)
				userHandler := handlers.NewUserRoleHandler(userRepo)
				adminGroup.GET("/AllUsersList", userHandler.ListAllUsers)
				adminGroup.GET("/UpdateRole", userHandler.UpdateRole)
			}
		}

	}

	return r
}

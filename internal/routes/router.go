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

	// Static Public Route
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "TendaBOX",
		})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "TendaBOX| Dashboard",
		})
	})
	r.GET("/sample", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sample.html", nil)
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "TendaBOX| Registeration",
		})
	})

	//404

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

			userRepo := repositroy.NewUserRepository(database.DB)
			profileHandler := &handlers.ProfileHandler{Repo: userRepo}

			// آدرس نهایی: /api/v1/user/myprofile
			userGroup.GET("/myprofile", profileHandler.GetProfile)

			// روت HTML پروفایل (اگر نامش را عوض کنید بهتر است تا با API تداخل نکند)
			userGroup.GET("/view-profile", func(c *gin.Context) {
				c.HTML(200, "UserProfile.html", gin.H{"title": "TendaBOX| Profile"})
			})

			//Admin Routes
			adminGroup := userGroup.Group("/admin")
			adminGroup.Use(middleware.AuthorizeRole("Admin", "super_admin"))
			{
				adminGroup.GET("/security", func(c *gin.Context) {
					c.HTML(200, "admin_ChangeUsersRole.html", gin.H{
						"title": "Manage User's Roles",
					})

				})
				adminGroup.GET("/usermangment", func(c *gin.Context) {
					c.HTML(200, "admin_ActiveUsers.html", gin.H{
						"title": "User Managment",
					})

				})
				userRepo := repositroy.NewUserRepository(database.DB)
				userHandler := handlers.NewUserRoleHandler(userRepo)
				adminGroup.GET("/AllUsersList", userHandler.ListAllUsers)
				adminGroup.POST("/UpdateRole", userHandler.UpdateRole)

				updatestatus_handler := handlers.NewUpdateUsersStatus(userRepo)
				adminGroup.POST("/ToggleActivation", updatestatus_handler.ToggleActivation)
			}
		}

	}

	return r
}

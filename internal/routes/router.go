package routes

import (
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
		c.HTML(http.StatusOK, "dashboard.html", nil)
	})

	//v1 := r.Group("/api/v1")
	//v1.POST("/login", )

	return r
}

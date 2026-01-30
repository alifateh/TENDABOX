package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAccessLevel(c *gin.Context) {
	userID, _ := c.Get("userID")
	userLevel, _ := c.Get("userLevel")
	slog.Info("Dashboard data requested", "user_id", userID, "User_Level", userLevel, "Client IP =", c.ClientIP())
	c.JSON(http.StatusOK, gin.H{
		"status": "authenticated",
		"level":  userLevel,
	})
}

package handlers

import (
	repo "Tendabox/internal/repository"
	"Tendabox/pkg/database"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GenerateMenu(c *gin.Context) {
	// گرفتن Role از context (که قبلاً توسط Middleware از JWT استخراج شده)
	roleName, exists := c.Get("userLevel")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
		return
	}

	// فراخوانی ریپازیتوری
	mainMenu, err := repo.GenrateMenu(roleName.(string), database.DB)
	if err != nil {
		slog.Error("Failed to generate menu", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch menu"})
		return
	}

	// ارسال منو به فرانت‌اند
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   mainMenu,
	})
}

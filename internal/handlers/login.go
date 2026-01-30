package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	repo "Tendabox/internal/repository"
	"Tendabox/pkg/auth"
	"Tendabox/pkg/database"

	"github.com/gin-gonic/gin"
)

type login_input struct {
	Email    string `binding:"required,email" json:"email" `
	Password string `binding:"required" json:"password"`
}

func Login(c *gin.Context) {
	var input login_input
	userRepo := repo.NewUserRepository(database.DB)
	if err := c.ShouldBindJSON(&input); err != nil {
		slog.Warn("Validation error", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email address is not well formatted"})
		return
	}

	email := strings.TrimSpace(input.Email)

	user, err := userRepo.GetByEmail(email)
	if err != nil {
		slog.Error("User not found", "email", email)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email Address is Wrong or Not Registered"})
		return
	}

	if err := user.CheckPassword(input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Password is Wrong"})
		return
	}

	roleName := "User"
	if user.Role.Name != "" {
		roleName = user.Role.Name
	}

	token, err := auth.GenerateToken(user.ID, roleName)
	if err != nil {
		slog.Error("JWT Error", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Craeting Token is Failed"})
		return
	}

	c.SetCookie("auth_token", token, 3600, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "success", "redirect": "/dashboard"})
}

package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type login_input struct {
	Email    string `binding:"required;email" json:"email" `
	Password string `bindingg:"required" json:"password"`
}

func Login(c *gin.Context) {
	var input login_input

	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		slog.Warn("Login attempt with invalid input", "error", err, "Client IP =", c.ClientIP())
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid Input"})
		return
	}
	email := strings.TrimSpace(input.Email)

	user, err := repo.userRepo.GetByEmail(email)

	if err != nil {
		slog.Error("User fetch failed", "email", email, "error", err, "Client IP =", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	if err := user.CheckPassword(input.Password); err != nil {
		slog.Warn("Incorrect password attempt", "email", email, "Client IP =", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Level)
	if err != nil {
		slog.Error("Failed to generate token", "user_id", user.ID, "err", err, "Client IP =", c.ClientIP())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token creation failed"})
		return
	}

	slog.Info("User logged in successfully", "user_id", user.ID, "level", user.Level, "Client IP =", c.ClientIP())

	//Cookie Set

	c.SetCookie("auth_token", token, 3600, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "success", "redirect": "/dashboard"})
}

package handlers

import (
	"Tendabox/internal/models"
	"Tendabox/internal/repositroy"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	repo repositroy.UserRepository
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{
		repo: repositroy.NewUserRepository(db),
	}
}

// این متد دقیقا امضای مورد نیاز Gin را دارد
func (h *UserHandler) RegisterUser(c *gin.Context) {
	val, _ := c.Get("validatedInput")

	// حالا این خط دیگر Panic نمی‌کند
	input := val.(models.RegisterInput)

	newUser := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password,
		RoleUUID:  input.RoleUUID,
	}

	if err := h.repo.CreateUser(&newUser); err != nil {
		if err.Error() == "user_already_exists" {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User Already Registered!!!",
			})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Registration successful wait until validation result"})
}

package handlers

import (
	"Tendabox/internal/models" // اصلاح غلط املایی پکیج (repository)
	repositroy "Tendabox/internal/repository"

	"github.com/gin-gonic/gin"
)

// تغییر نام برای جلوگیری از تداخل (Redeclaration)
type ProfileHandler struct {
	Repo repositroy.UserRepository
}

// سازنده برای هندلر (اختیاری اما استاندارد)
func NewProfileHandler(r repositroy.UserRepository) *ProfileHandler {
	return &ProfileHandler{Repo: r}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	var input models.UserProfile
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"message": "Data is NOT Well Formed!"})
		return
	}

	profile, err := h.Repo.GetUserProfile(input.UserID)
	if err != nil {
		c.JSON(500, gin.H{"error": "User Profile not found"})
		return
	}

	c.JSON(200, profile)
}

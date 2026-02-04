package handlers

import (
	middleware "Tendabox/internal/middelwars"
	repositroy "Tendabox/internal/repository"

	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	Repo repositroy.UserRepository
}

func NewProfileHandler(r repositroy.UserRepository) *ProfileHandler {
	return &ProfileHandler{Repo: r}
}

func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userPtr := middleware.CleanID(c)
	if userPtr == nil {
		return
	}

	profile, err := h.Repo.GetUserProfile(userPtr.ID)
	if err != nil {
		c.JSON(404, gin.H{"error": "Profile not found in database"})
		return
	}

	c.JSON(200, profile)
}

package handlers

import (
	"Tendabox/internal/models"
	repositroy "Tendabox/internal/repository"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type UpdateUserStatusHandler struct {
	reposit repositroy.UserRepository
}

func NewUpdateUsersStatus(r repositroy.UserRepository) *UpdateUserStatusHandler {
	return &UpdateUserStatusHandler{reposit: r}
}

func (u *UpdateUserStatusHandler) ToggleActivation(c *gin.Context) {
	var inputstatus models.UpdateUserStatus
	if err := c.ShouldBindJSON(&inputstatus); err != nil {
		// Log the raw error to see exactly why binding failed
		slog.Error("Binding failed", "error", err)
		c.JSON(400, gin.H{"message": "Error: Data is NOT Well Formed!"})
		return
	}

	// Now inputstatus.Active is already a bool (true/false)
	err_repo := u.reposit.UpdateUserStatus(inputstatus.UserID, inputstatus.Active)
	if err_repo != nil {
		c.JSON(500, gin.H{"message": "Update failed in database"})
		return
	}

	c.JSON(200, gin.H{"message": "Status updated successfully"})
}

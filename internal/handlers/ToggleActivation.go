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
		slog.Info("userID =" + c.GetString("userID"))
		slog.Info("active =" + c.GetString("active"))
		c.JSON(400, gin.H{
			"message": "Error: Data is NOT Well Formed!",
		})
		return
	}

	NewStatus := false
	if inputstatus.Active == "activate" {
		NewStatus = true
	}

	err_repo := u.reposit.UpdateUserStatus(inputstatus.UserID, NewStatus)
	if err_repo != nil {
		c.JSON(400, gin.H{
			"message": "Error: Data is NOT Well Formed!",
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "User status updated successfully",
	})
}

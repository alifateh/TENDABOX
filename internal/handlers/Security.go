package handlers

import (
	repositroy "Tendabox/internal/repository"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UpdateRoleInput struct {
	UserID string `json:"user_id" binding:"required,uuid"`
	RoleID string `json:"role_id" binding:"required,uuid"`
}

type UserRoleHandler struct {
	repo repositroy.UserRepository
}

func NewUserRoleHandler(r repositroy.UserRepository) *UserRoleHandler {
	return &UserRoleHandler{repo: r}
}

func (h *UserRoleHandler) UpdateRole(c *gin.Context) {
	var input UpdateRoleInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INPUT data is NOT Wellform"})
		return
	}

	err := h.repo.UpdateUserRole(input.UserID, input.RoleID)

	if err != nil {

		if err.Error() == "user not found" {
			slog.Error("User NOT Found!")
			c.JSON(http.StatusNotFound, gin.H{"error": "User Not Found"})
			return
		}
		slog.Error("Error updating User Role")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating User Role"})
		return
	}
	adminID, exists := c.Get("UserID")
	if !exists {
		adminID = "unknown"
	}

	slog.Info("User's Role changed Successfully",
		"Admin_IP", c.ClientIP(),
		"Admin_ID", adminID,
		"Target_User_ID", input.UserID,
	)
	c.JSON(http.StatusOK, gin.H{"message": "User's Role changed Successfully"})
}

func (h *UserRoleHandler) ListAllUsers(c *gin.Context) {
	list, err := h.repo.GetAllUsers()
	if err != nil {
		c.JSON(500, gin.H{"MSG": "Error in Users List"})
		return
	}
	c.JSON(200, list)
}

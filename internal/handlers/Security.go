package handlers

import (
	"Tendabox/internal/models"
	repositroy "Tendabox/internal/repository"
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserRoleHandler struct {
	repo repositroy.UserRepository
}

func NewUserRoleHandler(r repositroy.UserRepository) *UserRoleHandler {
	return &UserRoleHandler{repo: r}
}

func (h *UserRoleHandler) UpdateRole(c *gin.Context) {
	var input models.UpdateRoleInput
	slog.Info("userID =" + c.GetString("userID"))
	slog.Info("UserRole =" + c.GetString("UserRole"))

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
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
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	list, total, err := h.repo.GetAllUsers(page, limit)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error in Users List"})
		return
	}

	c.JSON(200, gin.H{
		"data":      list,
		"total":     total,
		"page":      page,
		"last_page": math.Ceil(float64(total) / float64(limit)),
	})
}

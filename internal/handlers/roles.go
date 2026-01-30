package handlers

import (
	"Tendabox/internal/models"
	repo "Tendabox/internal/repository"
	"Tendabox/pkg/database"

	"github.com/gin-gonic/gin"
)

func GetAllRoles(c *gin.Context) {
	var Obj_Roles []models.Roles
	Obj_Roles, err := repo.AllRoles(database.DB)
	if err != nil {
		c.JSON(500, gin.H{"Error": "Fatal Error"})
		return
	}
	c.JSON(200, Obj_Roles)

}

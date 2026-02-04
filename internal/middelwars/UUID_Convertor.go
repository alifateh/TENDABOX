package middleware

import (
	"Tendabox/internal/models"
	"log/slog"

	"github.com/gin-gonic/gin"
)

//fmt.Printf("Keys in context: %+v\n", c.Keys)
///Keys in context: map[userID:8b10089b-270e-4a62-b922-6f87deb554a4 userLevel:super_admin userROleID:b8616edd-624b-4257-b2f9-1f37dee7df58]

func CleanID(c *gin.Context) *models.User {

	IDval, idexists := c.Get("userID")
	roleIDval, roleexists := c.Get("userROleID")

	if !idexists {
		slog.Error("User ID missing in Gin Context ---> User ID")
		c.JSON(400, gin.H{"error": "your token is not valid!!! User ID"})
		return nil
	}

	if !roleexists {
		slog.Error("User Role ID was Not Exist in Token ---> Role ID")
		c.JSON(400, gin.H{"error": "your token is not valid!!! Role ID"})
		return nil
	}

	///role  userLevel:super_admin
	return &models.User{
		ID:       IDval.(string),
		RoleUUID: roleIDval.(string),
	}
}

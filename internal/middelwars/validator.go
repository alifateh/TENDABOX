package middleware

import (
	"Tendabox/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterValidator(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Input Validation Failed!"})
		c.Abort()
		return
	}
	c.Set("validatedInput", input)
	c.Next()
}

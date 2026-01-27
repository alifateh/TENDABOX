package middleware

import (
	"Tendabox/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterValidator() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input models.RegisterInput // از مدل نام‌دار استفاده کنید

		if err := c.ShouldBindJSON(&input); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "Validation failed: " + err.Error()})
			return
		}

		c.Set("validatedInput", input)
		c.Next()
	}
}

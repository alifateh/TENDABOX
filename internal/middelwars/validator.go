package middleware

import (
	"Tendabox/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func RegisterValidator(c *gin.Context) {
	var input models.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {

		if errs, ok := err.(validator.ValidationErrors); ok {
			listErrors := make(map[string]string)

			for _, e := range errs {
				listErrors[e.Field()] = "Error in tag :" + e.Tag()
			}

			c.JSON(http.StatusBadRequest, gin.H{"errors": listErrors})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Input Wrong JSON Format"})
		}

		c.Abort()
		return
	}

	c.Set("validatedInput", input)
	c.Next()
}

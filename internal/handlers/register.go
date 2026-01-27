package handlers

import (
	"Tendabox/internal/models"

	"gorm.io/gorm"
)

func RegisterUser(*gorm.DB) (bool, err) {
	newUser := models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Password:  input.Password, // نکته: حتما پسورد را Hash کنید (مثلا با bcrypt)
		RoleUUID:  input.RoleUUID,
	}

}

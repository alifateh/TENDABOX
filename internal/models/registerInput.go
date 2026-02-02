package models

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required,min=3,alphaunicode"`

	LastName string `json:"last_name" binding:"required,min=3,alphaunicode"`

	Email string `json:"email" binding:"required,email"`

	Password string `json:"password" binding:"required,min=6,alphanum"`

	RoleUUID string `json:"role_uuid" binding:"required,uuid4"`
}

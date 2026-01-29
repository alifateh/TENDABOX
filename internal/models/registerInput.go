package models

type RegisterInput struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"-" binding:"required,min=6"`
	RoleUUID  string `json:"role_uuid" binding:"required"`
}

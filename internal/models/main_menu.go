package models

import "gorm.io/gorm"

type MainMenu struct {
	gorm.Model
	ItemName string `gorm:"column:item_name;type:varchar(100)" json:"item_name" binding:"required,min=3,max=25"`
	URLPath  string `gorm:"column:url;type:varchar(255)" json:"url" binding:"required"`
	MinLevel int    `gorm:"column:min_level" json:"min_level" binding:"required"`
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Password string `gorm:"type:varchar(360)"json:"password"`
}

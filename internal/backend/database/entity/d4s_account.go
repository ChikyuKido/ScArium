package entity

import "gorm.io/gorm"

type D4sAccount struct {
	gorm.Model
	UserId      uint
	Username    string
	Password    string `json:"-"`
	DisplayName string
	ImageUrl    string
}

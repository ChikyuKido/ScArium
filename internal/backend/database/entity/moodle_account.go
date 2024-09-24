package entity

import "gorm.io/gorm"

type MoodleAccount struct {
	gorm.Model
	UserId      uint
	InstanceUrl string
	Username    string
	Password    string `json:"-"`
	DisplayName string
	ImageUrl    string
}

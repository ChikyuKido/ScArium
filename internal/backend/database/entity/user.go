package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string
	Password string `gorm:"size:255" json:"-"`
	Admin    bool
}

package entity

type User struct {
	Model
	Username string
	Password string `gorm:"size:255" json:"-"`
	Admin    bool
}

package entity

type User struct {
	Model
	Username       string
	Password       string `gorm:"size:255" json:"-"`
	Admin          bool
	MoodleAccounts []MoodleAccount `gorm:"foreignKey:UserID"`
	D4sAccount     []D4sAccount    `gorm:"foreignKey:UserID"`
}

package repo

import (
	"ScArium/internal/backend/database"
	"ScArium/internal/backend/database/entity"
)

func GetUserByName(username string) (*entity.User, error) {
	var u entity.User
	if err := database.GetDB().Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func DoesUserByNameExists(username string) bool {
	var count int64
	database.GetDB().Model(&entity.User{}).Where("username = ?", username).Count(&count)
	return count > 0
}

func InsertNewUser(user entity.User) error {
	if err := database.GetDB().Create(&user).Error; err != nil {
		return err
	}
	return nil
}

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

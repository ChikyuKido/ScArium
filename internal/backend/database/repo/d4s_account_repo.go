package repo

import (
	"ScArium/internal/backend/database"
	"ScArium/internal/backend/database/entity"
)

func CreateD4sAccount(user *entity.User, username string, password string, displayName string, imageId string) error {
	d4s := entity.D4sAccount{
		UserID:      user.ID,
		Username:    username,
		Password:    password,
		DisplayName: displayName,
		ImageId:     imageId,
	}

	return database.GetDB().Create(&d4s).Error
}

func GetD4sAccounts(user *entity.User) ([]entity.D4sAccount, error) {
	var accounts []entity.D4sAccount
	if err := database.GetDB().Where(&entity.D4sAccount{UserID: user.ID}).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

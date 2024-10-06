package repo

import (
	"ScArium/internal/backend/database"
	"ScArium/internal/backend/database/entity"
)

func CreateMoodleAccount(user *entity.User, instanceUrl string, username string, password string, displayName string, imageId string) error {
	moodle := entity.MoodleAccount{
		UserID:      user.ID,
		InstanceUrl: instanceUrl,
		Username:    username,
		Password:    password,
		DisplayName: displayName,
		ImageId:     imageId,
	}

	return database.GetDB().Create(&moodle).Error
}
func GetMoodleAccountById(id uint) (*entity.MoodleAccount, error) {
	moodle := entity.MoodleAccount{}
	if err := database.GetDB().Where("id = ?", id).First(&moodle).Error; err != nil {
		return nil, err
	}
	return &moodle, nil
}
func GetMoodleAccounts(user *entity.User) ([]entity.MoodleAccount, error) {
	var moodleAccounts []entity.MoodleAccount
	if err := database.GetDB().Where(&entity.MoodleAccount{UserID: user.ID}).Find(&moodleAccounts).Error; err != nil {
		return nil, err
	}
	return moodleAccounts, nil
}

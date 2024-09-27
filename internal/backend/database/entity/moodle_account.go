package entity

type MoodleAccount struct {
	Model
	UserId      uint   `json:"user_id"`
	InstanceUrl string `json:"instance_url"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	DisplayName string `json:"display_name"`
	ImageId     string `json:"image_id"`
}

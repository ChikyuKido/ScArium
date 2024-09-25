package entity

type D4sAccount struct {
	Model
	UserId      uint   `json:"user_id"`
	Username    string `json:"username"`
	Password    string `json:"-"`
	DisplayName string `json:"display_name"`
	ImageUrl    string `json:"image_url"`
}

package config

import (
	"ScArium/internal/backend/helper"
	"os"
)

var (
	DATA_PATH            = "./data"
	DATABASE_PATH        = DATA_PATH + "/database.db"
	STATIC_PATH          = "./static"
	ADMIN_REGISTER_EXIST = DATA_PATH + "/adminRegister.exists"
	IMAGE_PATH           = DATA_PATH + "/images"
	JOBS_PATH            = DATA_PATH + "/jobs"
	ACCOUNTS_PATH        = DATA_PATH + "/accounts"
	MOODLE_PATH          = ACCOUNTS_PATH + "/moodle"
)
var (
	RT_ADMIN_REGISTER_AVAILABLE = false
)

func InitConfig() {
	os.MkdirAll(DATA_PATH, os.ModePerm)
	os.MkdirAll(IMAGE_PATH, os.ModePerm)
	os.MkdirAll(JOBS_PATH, os.ModePerm)
	os.MkdirAll(ACCOUNTS_PATH, os.ModePerm)
	os.MkdirAll(MOODLE_PATH, os.ModePerm)
	RT_ADMIN_REGISTER_AVAILABLE = !helper.DoesFileExists(ADMIN_REGISTER_EXIST)
}

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
)
var (
	RT_ADMIN_REGISTER_AVAILABLE = false
)

func InitConfig() {
	os.MkdirAll(DATA_PATH, os.ModePerm)
	RT_ADMIN_REGISTER_AVAILABLE = !helper.DoesFileExists(ADMIN_REGISTER_EXIST)
}

package config

import "os"

var (
	DATA_PATH     = "./data"
	DATABASE_PATH = DATA_PATH + "/database.db"
	STATIC_PATH   = "./static"
)

func InitConfig() {
	os.MkdirAll(DATA_PATH, os.ModePerm)
	os.MkdirAll(STATIC_PATH, os.ModePerm)
}

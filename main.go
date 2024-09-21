package main

import (
	"ScArium/common/config"
	"ScArium/common/log"
	"ScArium/internal"
)

func main() {
	config.InitConfig()
	log.InitLogger()
	server := internal.NewServer(7665, "localhost")
	server.Start()
}

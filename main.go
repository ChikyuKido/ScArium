package main

import (
	"ScArium/common/config"
	"ScArium/common/log"
	"ScArium/internal/backend/database"
	"ScArium/internal/jobs"
	"ScArium/internal/jobs/moodejob"
	"fmt"
	"time"
)

func sleepy() {
	for {
		time.Sleep(5 * time.Second)
	}
}

func main() {
	config.InitConfig()
	log.InitLogger()
	go jobs.Worker()
	database.InitDB()
	//server := internal.NewServer(7665, "localhost")
	//server.Start()
	err := moodejob.NewMoodleSyncJob(1)
	if err != nil {
		fmt.Println(err)
	}
	sleepy()
}

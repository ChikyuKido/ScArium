package main

import (
	"ScArium/common/config"
	"ScArium/common/log"
	"ScArium/external/moodle/mFunctions"
)

func main() {
	config.InitConfig()
	log.InitLogger()

	courses, err := mFunctions.GetCourses(mc)
	if err != nil {
		return
	}
	mFunctions.GetSections(mc, courses[0])

	//server := internal.NewServer(7665, "localhost")
	//server.Start()
}

package moodejob

import (
	"ScArium/external/moodle"
	"ScArium/external/moodle/mModel"
	"ScArium/internal/jobs"
	"fmt"
	"github.com/sirupsen/logrus"
)

func syncFiles(job jobs.Job, courses []*mModel.Course, id uint64, mc *moodle.MoodleClient) {
	sections := 0
	modules := 0
	for _, course := range courses {
		for _, section := range course.Sections {
			modules += len(section.Modules)
		}
		sections += len(course.Sections)
	}
	job.AppendLog(logrus.InfoLevel, fmt.Sprintf("Found %d courses with %d sections and %d modules", len(courses), sections, modules))
	job.AppendLog(logrus.InfoLevel, "Start syncing them")

	//basePath := fmt.Sprintf("%s/%d", config.MOODLE_PATH, id)
	//
	//for i, course := range courses {
	//
	//}

}

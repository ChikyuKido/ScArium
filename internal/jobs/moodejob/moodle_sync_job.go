package moodejob

import (
	"ScArium/common/config"
	"ScArium/external/moodle"
	"ScArium/external/moodle/mFunctions"
	"ScArium/external/moodle/mModel"
	"ScArium/internal/backend/database/repo"
	"ScArium/internal/jobs"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

func NewMoodleSyncJob(accountId uint) error {
	err := jobs.CreateJob("MoodleSync", "For accountId "+strconv.Itoa(int(accountId)),
		map[string]string{
			"accountId": strconv.Itoa(int(accountId)),
		},
		syncMoodle,
	)
	if err != nil {
		return fmt.Errorf("failed to create moodle sync job: %w", err)
	}
	return nil
}

func syncMoodle(job jobs.Job) error {
	accountId, err := strconv.ParseUint(job.Params["accountId"], 10, 32)
	if err != nil {
		return fmt.Errorf("failed to convert account id to uint: %w", err)
	}
	moodleAcc, err := repo.GetMoodleAccountById(uint(accountId))
	if err != nil {
		return fmt.Errorf("failed to get moodle account: %w", err)
	}
	mc, err := moodle.NewMoodleClient(moodleAcc.InstanceUrl, moodleAcc.Username, moodleAcc.Password)
	if err != nil {
		return fmt.Errorf("failed to create moodle client: %v", err)
	}
	courses, err := mFunctions.GetCourses(mc)
	if err != nil {
		return fmt.Errorf("failed to get courses: %v", err)
	}
	job.AppendLog(logrus.InfoLevel, fmt.Sprintf("Found %d courses", len(courses)))
	replaceImagesForCourses(job, courses)
	retrieveSections(job, courses, mc)
	createMetadata(job, courses, accountId)
	return nil
}

func createMetadata(job jobs.Job, courses []*mModel.Course, accountId uint64) {
	basePath := config.MOODLE_PATH + "/" + strconv.FormatUint(accountId, 10)
	err := os.MkdirAll(basePath, os.ModePerm)
	if err != nil {
		job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata directory: %v", err))
		return
	}
	jsonData, err := json.Marshal(courses)
	if err != nil {
		job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to marshal courses: %v", err))
		return
	}
	err = os.WriteFile(basePath+"/courses.json", jsonData, 0666)
	if err != nil {
		job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata file: %v", err))
		return
	}
	for _, course := range courses {
		coursePath := basePath + "/" + strconv.Itoa(course.ID)
		err := os.MkdirAll(coursePath, os.ModePerm)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata directory: %v", err))
			return
		}
		for _, section := range course.Sections {
			jsonData, err := json.Marshal(section)
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to marshal courses: %v", err))
				return
			}
			err = os.WriteFile(coursePath+"/section_"+strconv.Itoa(section.ID)+".json", jsonData, 0666)
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("Failed to create metadata file: %v", err))
				return
			}
		}
	}
}

func retrieveSections(job jobs.Job, courses []*mModel.Course, mc *moodle.MoodleClient) {
	for _, course := range courses {
		err := mFunctions.GetSections(mc, course)
		if err != nil {
			job.AppendLog(logrus.ErrorLevel, fmt.Sprintf("failed to get sections for course: %s: %v", course.ShortName, err))
		}
		job.AppendLog(logrus.InfoLevel, fmt.Sprintf("Found sections %d for course: %s", len(course.Sections), course.ShortName))
	}
}

func replaceImagesForCourses(job jobs.Job, courses []*mModel.Course) {
	for i := range courses {
		if strings.Contains(courses[i].CourseImage, "http") {
			file, err := os.ReadFile("static/imgs/defaultMoodleCourse.svg")
			if err != nil {
				job.AppendLog(logrus.ErrorLevel, "Could not load default svg")
			}
			job.AppendLog(logrus.WarnLevel, fmt.Sprintf("%s does not have a image load default", courses[i].ShortName))
			courses[i].CourseImage = string(file)
		}
	}
}
